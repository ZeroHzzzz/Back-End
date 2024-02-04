package studentcontroller

import (
	"hr/app/service"
	"hr/app/utils"
	"hr/configs/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetConcreteSorceHandler(c *gin.Context) {
	// 上传申报
	c.Header("Content-Type", "application/json")
	userId := c.Param("userId")
	academicYear := c.Param("academicYear")

	// 从上下文中获取mongo客户端
	filter := bson.M{
		"userId":       userId,
		"academicTear": academicYear,
	}

	var result models.Score
	err := service.FindOne(c, "", "", filter).Decode(&result)
	if err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}
	utils.ResponseSuccess(c, result)
}

func GetYearScoreHandler(c *gin.Context) {
	// 传入userId，在student库中查找出对应的学生信息，返回map[string]int类型grade
	// 上传申报
	c.Header("Content-Type", "application/json")
	userId := c.Param("userId")

	filter := bson.M{
		"userId": userId,
	}

	var student models.Student
	err := service.FindOne(c, "", "", filter).Decode(&student)
	if err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}
	utils.ResponseSuccess(c, student.Mark)
}
