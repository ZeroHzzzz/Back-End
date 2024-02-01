package counsellorcontroller

import (
	"context"
	"hr/app/utils"
	"hr/configs/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	const DatabaseName string = ""
	const CollectionName string = ""

	var getsubmissionlistinformation getSubmissionListInformation
	err := c.ShouldBindJSON(&getsubmissionlistinformation)
	if err != nil {
		utils.ResponseError(c, "Paramter", "ParameterErrorMsg")
		return
	}

	// 获取collection
	var list []models.SubmitInformation
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	database := mongoClient.Database(DatabaseName)
	collection := database.Collection(CollectionName)
	// 这里可能会有bug，因为这里是嵌套字段，不知道能不能直接查出来
	// 获取未审核表单
	filter := bson.M{
		"class":      getsubmissionlistinformation.Class,
		"profession": getsubmissionlistinformation.Profession,
		"grade":      getsubmissionlistinformation.Grade,
		"status":     false,
	}
	options := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip((getsubmissionlistinformation.Index - 1) * getsubmissionlistinformation.PaginationSize).SetLimit(getsubmissionlistinformation.PaginationSize)

	// 执行查询
	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		// 处理
		return
	}
	if err := cursor.All(context.TODO(), &list); err != nil {
		return
	}

	utils.ResponseSuccess(c, list)
	return
}

func GetSubmission(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	const DatabaseName string = ""
	const CollectionName string = "" //student

	submissionId := c.Param("submissionId")

	// 从上下文中获取mongo客户端
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	database := mongoClient.Database("DatabaseName")
	collection := database.Collection("CollectionName")
	filter := bson.M{"submissionId": submissionId}
	result := collection.FindOne(c, filter)

	var submission models.SubmitInformation
	err := result.Decode(&submission)
	if err != nil {
		// 处理
		return
	}
	utils.ResponseSuccess(c, submission)
}
