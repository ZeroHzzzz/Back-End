package counsellorhandler

import (
	"context"
	"hr/app/service"
	"hr/app/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type AccessTime struct {
	StartTime string `json:"StartTime"`
	EndTime   string `json:"EndTime"`
}

// 管理员设置访问时间段的处理函数
func SetAccessTimeHandler(c *gin.Context) {
	var accessTime AccessTime
	if err := c.BindJSON(&accessTime); err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}
	startTime, err := time.Parse("2006-01-02 15:04", accessTime.StartTime)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}
	endTime, err := time.Parse("2006-01-02 15:04", accessTime.EndTime)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}
	redisClient := service.GetRedisClint(c)
	redisClient.Set(context.Background(), "Start-Time", startTime, 0)
	redisClient.Set(context.Background(), "End-Time", endTime, 0)
	utils.ResponseSuccess(c, nil)
}
