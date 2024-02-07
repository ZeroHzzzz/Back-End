package midware

import (
	"errors"
	"hr/app/utils"
	"hr/configs/models"
	"math/rand"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	models.CurrentUser
	jwt.StandardClaims
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const TokenExpireDuration = time.Minute * 15 // 设置过期时间

// 生成随机密钥
func generateRandomString(length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rng.Intn(len(charset))]
	}

	return string(result)
}

var jwtKey = []byte(generateRandomString(10))

// 生成token
func GenerateToken(currentUser models.CurrentUser) (string, error) {
	expirationTime := time.Now().Add(30 * time.Minute) // 新Token有效期为30分钟

	claims := &Claims{
		CurrentUser: currentUser,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "ZeroHzzzz", //签发人，可以修改
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*Claims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func JWTAuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("authorization")
		if authHeader == "" {
			utils.ResponseUnauthorized(c)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.Split(authHeader, ".")
		if len(parts) != 3 {
			utils.ResponseUnauthorized(c)
			c.Abort()
			return
		}
		mc, err := ParseToken(string(jwtKey))
		if err != nil {
			utils.ResponseUnauthorized(c)
			c.Abort()
			return
		}

		// 检查用户的角色是否在允许的角色列表中
		roleAllowed := false
		for _, allowedRole := range allowedRoles {
			if mc.CurrentUser.Role == allowedRole {
				roleAllowed = true
				break
			}
		}
		if !roleAllowed {
			utils.ResponseUnauthorized(c)
			c.Abort()
			return
		}

		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("CurrentUser", mc.CurrentUser)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
