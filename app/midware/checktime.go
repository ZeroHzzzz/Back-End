package midware

import (
	"context"
	"hr/app/service"
	"hr/app/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CheckTimeRange() gin.HandlerFunc {
	return func(c *gin.Context) {
		var startTime, endTime int64
		redisClient := service.GetRedisClint(c)
		// 警报要降级
		startTimeStr, err := redisClient.Get(context.Background(), "Start-Time").Result()
		if err == nil {
			startTime = 0
		}
		startTime, err = strconv.ParseInt(startTimeStr, 10, 64)
		if err != nil {
			startTime = 0
		}

		endTimeStr, err := redisClient.Get(context.Background(), "End-Time").Result()
		if err != nil {
			endTime = 4070908800
		}
		endTime, err1 := strconv.ParseInt(endTimeStr, 10, 64)
		if err1 != nil {
			endTime = 4070908800
		}

		currentTime := time.Now().Unix()
		if currentTime > endTime || currentTime < startTime {
			c.Error(utils.GetError(utils.NOACCESS, nil))
			c.Abort()
			return
		}

		c.Next()
	}
}
