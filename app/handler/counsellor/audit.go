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

type auditOneInformation struct {
	Status bool   `json:"status"`
	Cause  string `json:"cause"`
	Advice string `json:"advice"`
}

func AuditOne(c *gin.Context) {
	// 审批单个申报
	c.Header("Content-Type", "application/json")

	// 获取用户信息
	currentUser := service.GetCurrentUser(c)
	const DatabaseName string = ""
	const CollectionName string = "" //submission

	var information auditOneInformation
	err := c.ShouldBindJSON(&information)
	if err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}
	submissionId := c.Param("submissionId")

	// 从上下文中获取mongo客户端
	filter := bson.M{
		"submissionId": submissionId,
	}
	modified := bson.M{
		"$set": bson.M{
			"status": information.Status,
			"cause":  information.Cause,
			"advice": information.Advice,
		},
	}
	_ = service.UpdateOne(c, DatabaseName, CollectionName, filter, modified)
	// 新建历史记录
	newHistory := models.SubmitHistory{
		SubmissionId: submissionId,
		AuditorId:    currentUser.UserId,
		Status:       information.Status,
		Cause:        information.Cause,
		Advice:       information.Advice,
	}
	_ = service.InsertOne(c, DatabaseName, CollectionName, newHistory)
	utils.ResponseSuccess(c, nil)
	return
}

type auditManyInformation struct {
	SubmissionIds []string `json:"submissionIds"`
	Status        bool     `json:"status"`
	Advice        string   `json:"advice"`
	Cause         string   `json:"cause"`
}

func AuditManySubmission(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	// 获取用户信息
	currentUser := service.GetCurrentUser(c)
	const DatabaseName string = ""
	const CollectionName string = "" //submission
	var information auditManyInformation
	err := c.ShouldBindJSON(&information)
	if err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}
	// 从上下文中获取mongo客户端
	filter := bson.M{
		"_id": bson.M{
			"$in": information.SubmissionIds,
		},
	}
	modified := bson.M{
		"$set": bson.M{
			"status": information.Status,
			"cause":  information.Cause,
			"advice": information.Advice,
		},
	}
	_ = service.UpdateMany(c, DatabaseName, CollectionName, filter, modified)
	// 记录历史申报
	baseSubmission := models.SubmitHistory{
		AuditorId: currentUser.UserId,
		Status:    information.Status,
		Cause:     information.Cause,
		Advice:    information.Advice,
	}
	var submissions []interface{}
	for _, submissionId := range information.SubmissionIds {
		doc := baseSubmission
		doc.SubmissionId = submissionId
		submissions = append(submissions, doc)
	}

	_ = service.InsertMany(c, DatabaseName, CollectionName, submissions)
	utils.ResponseSuccess(c, nil)
}

type getAuditlist struct {
	Index          int64 `json:"index"`
	PaginationSize int64 `json:"paginationSize"`
}

func GetAuditHistory(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var information getAuditlist
	const DatabaseName string = ""
	const CollectionName string = "" //student

	if err := c.ShouldBindJSON(&information); err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}

	filter := bson.D{}
	options := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip((information.Index - 1) * information.PaginationSize).SetLimit(information.PaginationSize)
	result := service.Find(c, DatabaseName, CollectionName, filter, options)

	var list []models.SubmitHistory
	if err := result.All(context.TODO(), &list); err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}
	utils.ResponseSuccess(c, list)
	return
}
