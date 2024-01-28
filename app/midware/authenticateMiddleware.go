package midware

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 生成随机密钥
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Grade    int    `json:"grade"`
	jwt.StandardClaims
}

func AuthenticateMiddleware(c *gin.Context, allowedRoles ...string) {
	authorizationHeader := c.GetHeader("Authorization")
	if authorizationHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", 1)
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// 检查用户的角色是否在允许的角色列表中
	roleAllowed := false
	for _, allowedRole := range allowedRoles {
		if claims.Role == allowedRole {
			roleAllowed = true
			break
		}
	}

	if !roleAllowed {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// // 检查特定角色的管理范围
	// if claims.Role == "admin" && claims.Grade != claims.Grade {
	// 	c.AbortWithStatus(http.StatusForbidden)
	// 	return
	// }
}
