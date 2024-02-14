package studenthandler

import (
	"context"
	"hr/app/service"
	"hr/app/utils"
	"hr/configs/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetConcreteSorce(c *gin.Context) {
	// 上传申报
	c.Header("Content-Type", "application/json")
	userId := c.Query("userId")
	academicYear := c.Query("academicYear")
	// 从上下文中获取mongo客户端
	filter := bson.M{
		"userId":       userId,
		"academicYear": academicYear,
	}

	var score []models.Score
	result := service.Find(c, utils.MongodbName, utils.Score, filter)
	if err := result.All(context.TODO(), &score); err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return
	}
	utils.ResponseSuccess(c, score)
}

// func GetYearScoreHandler(c *gin.Context) {
// 	// 传入userId，在student库中查找出对应的学生信息，返回map[string]int类型grade
// 	// 上传申报
// 	c.Header("Content-Type", "application/json")
// 	userId := c.Param("userId")

// 	filter := bson.M{
// 		"userId": userId,
// 	}

// 	var student models.Student
// 	err := service.FindOne(c, utils.MongodbName, utils.Score, filter).Decode(&student)
// 	if err != nil {
// 		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
// 		return
// 	}
// 	utils.ResponseSuccess(c, student.Mark)
// }
