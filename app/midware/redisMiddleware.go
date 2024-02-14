package midware

import (
	"context"
	"hr/app/utils"
	configs "hr/configs/configuration"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func RedisClientMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		redisurl := configs.Config.GetString("redis.url")
		redispassword := configs.Config.GetString("redis.password")
		clientOptions := redis.Options{
			Addr:     redisurl,
			Password: redispassword,
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
