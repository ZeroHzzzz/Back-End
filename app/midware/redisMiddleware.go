package midware

import (
	"context"
	"fmt"
	"hr/app/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const (
	redisHost     = "localhost"
	redisPort     = 6379
	redisPassword = ""
)

func redisClientMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientOptions := redis.Options{
			Addr:     fmt.Sprintf("%s:%d", redisHost, redisPort),
			Password: redisPassword,
		}
		client := redis.NewClient(&clientOptions)
		defer client.Close()

		if err := client.Ping(context.Background()).Err(); err != nil {
			c.Error(utils.GetError(utils.CONNECT_ERROR, err.Error()))
			c.Abort()
			return
		}
		ctx := context.WithValue(c.Request.Context(), "redisClient", client)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
