package midware

import (
	"context"
	"hr/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func redisClientMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		client, err := service.InitRedis()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize redis client"})
			c.Abort()
			return
		}
		ctx := context.WithValue(c.Request.Context(), "redisClient", client)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
