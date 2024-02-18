package counsellorhandler

import (
	"context"
	"hr/app/service"
	"hr/app/utils"

	"github.com/gin-gonic/gin"
)

type accessTime struct {
	StartTime int64 `json:"StartTime"`
	EndTime   int64 `json:"EndTime"`
}

// 管理员设置访问时间段的处理函数
func SetAccessTimeHandler(c *gin.Context) {
	var accessTime accessTime
	if err := c.ShouldBindJSON(&accessTime); err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}
	redisClient := service.GetRedisClint(c)
	redisClient.Set(context.Background(), "Start-Time", accessTime.StartTime, 0)
	redisClient.Set(context.Background(), "End-Time", accessTime.EndTime, 0)
	utils.ResponseSuccess(c, nil)
}
