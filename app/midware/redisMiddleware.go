package midware

import (
	"context"
	"fmt"
	"net/http"

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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize MongoDB client"})
			c.Abort()
		}
		ctx := context.WithValue(c.Request.Context(), "redisClient", client)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
