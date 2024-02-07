package midware

import (
	"context"
	"hr/app/service"
	"hr/app/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func CheckTimeRange() gin.HandlerFunc {
	return func(c *gin.Context) {
		var startTime, endTime time.Time
		redisClient := service.GetRedisClint(c)
		startTime, err := redisClient.Get(context.Background(), "Start-Time").Time()
		if err == redis.Nil {
			startTime = time.Date(1999, time.January, 1, 1, 0, 0, 0, time.UTC)
		}
		endTime, err = redisClient.Get(context.Background(), "End-Time").Time()
		if err == redis.Nil {
			startTime = time.Date(3000, time.January, 1, 1, 0, 0, 0, time.UTC)
		}

		currentTime := time.Now()

		if currentTime.Before(startTime) || currentTime.After(endTime) {
			c.Error(utils.GetError(utils.NOACCESS, nil))
			c.Abort()
			return
		}

		c.Next()
	}
}
