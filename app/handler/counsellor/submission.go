package counsellor

import (
	"context"
	"hr/app/service"
	"hr/app/utils"
	"hr/configs/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type getSubmissionListInformation struct {
	Index          int64 `json:"index" binding:"required"`
	PaginationSize int64 `json:"paginationSize"`
	Profession     int64 `json:"profession"`
	Grade          int64 `json:"grade"` //年级
	Class          int64 `json:"class"`
}

func GetSubmissionList(c *gin.Context) {
	var getsubmissionlistinformation getSubmissionListInformation
	err := c.ShouldBindJSON(&getsubmissionlistinformation)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}

	// 获取collection
	var list []models.SubmitInformation

	// 获取未审核表单
	filter := bson.M{
		"class":      getsubmissionlistinformation.Class,
		"profession": getsubmissionlistinformation.Profession,
		"grade":      getsubmissionlistinformation.Grade,
		"status":     false,
	}
	options := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip((getsubmissionlistinformation.Index - 1) * getsubmissionlistinformation.PaginationSize).SetLimit(getsubmissionlistinformation.PaginationSize)

	// 执行查询
	cursor := service.Find(c, utils.MongodbName, utils.Submission, filter, options)
	if err := cursor.All(context.TODO(), &list); err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return
	}
	utils.ResponseSuccess(c, list)
	return
}

func GetSubmission(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	submissionId := c.Param("submissionId")

	filter := bson.M{"submissionId": submissionId}
	result := service.FindOne(c, utils.MongodbName, utils.Submission, filter)

	var submission models.SubmitInformation
	err := result.Decode(&submission)
	if err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return
	}
	utils.ResponseSuccess(c, submission)
}
