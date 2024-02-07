package counsellor

import (
	"context"
	"hr/app/service"
	"hr/app/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AccessTime struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// 管理员设置访问时间段的处理函数
func SetAccessTimeHandler(c *gin.Context) {
	var accessTime AccessTime
	if err := c.BindJSON(&accessTime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	startTime, err := time.Parse(time.RFC1123, accessTime.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format"})
		return
	}
	endTime, err := time.Parse(time.RFC1123, accessTime.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format"})
		return
	}
	redisClient := service.GetRedisClint(c)
	redisClient.Set(context.Background(), "Start-Time", startTime, 0)
	redisClient.Set(context.Background(), "End-Time", endTime, 0)
	utils.ResponseSuccess(c, nil)
}
