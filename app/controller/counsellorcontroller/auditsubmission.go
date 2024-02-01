package counsellorcontroller

import (
	"context"
	"fmt"
	"hr/app/utils"
	"hr/configs/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type auditOneSubmissionInformation struct {
	Status bool   `json:"status"`
	Cause  string `json:"cause"`
	Advice string `json:"advice"`
}

func AuditOneSubmission(c *gin.Context) {
	// 审批单个申报
	c.Header("Content-Type", "application/json")

	const DatabaseName string = ""
	const CollectionName string = "" //submission
	var information auditOneSubmissionInformation
	err := c.ShouldBindJSON(&information)
	if err != nil {
		utils.ResponseError(c, "failure", "Parameter wrong")
		return
	}
	submissionId := c.Param("submissionId")
	// 从上下文中获取mongo客户端
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	database := mongoClient.Database("DatabaseName")
	collection := database.Collection("CollectionName")
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
	_, err = collection.UpdateOne(context.TODO(), filter, modified)
	if err != nil {
		//处理逻辑
		if err == mongo.ErrNoDocuments {
			fmt.Println("No matching document found")
			return
		}
		log.Fatal(err)
		return
	}
	// 获取用户信息
	user, ok := c.Get("currentUser")
	currentUser, ok := user.(models.CurrentUser)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 新建历史记录
	newHistory := models.SubmitHistory{
		SubmissionId: submissionId,
		AuditorId:    currentUser.UserId,
		Status:       information.Status,
		Cause:        information.Cause,
		Advice:       information.Advice,
	}
	_, err = collection.InsertOne(context.Background(), newHistory)
	if err != nil {
		log.Fatal(err)
	}
	utils.ResponseSuccess(c, nil)
	return
}

type auditManySubmissionInformation struct {
	SubmissionIds []string `json:"submissionIds"`
	Status        bool     `json:"status"`
	Advice        string   `json:"advice"`
	Cause         string   `json:"cause"`
}

type auditManyResponse struct {
	SuccessCount int    `json:"successCount"`
	FailureCount int    `json:"failureCount"`
	Error        string `json:"error,omitempty"`
}

func AuditManySubmission(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	const DatabaseName string = ""
	const CollectionName string = "" //submission
	var information auditManySubmissionInformation
	err := c.ShouldBindJSON(&information)
	if err != nil {
		utils.ResponseError(c, "failure", "Parameter wrong")
		return
	}
	// 从上下文中获取mongo客户端
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	database := mongoClient.Database("DatabaseName")
	collection := database.Collection("CollectionName")
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
	_, err = collection.UpdateMany(context.TODO(), filter, modified)
	if err != nil {
		//处理逻辑
		if err == mongo.ErrNoDocuments {
			fmt.Println("No matching document found")
			return
		}
		log.Fatal(err)
		return
	}
	// 获取用户信息
	user, ok := c.Get("currentUser")
	currentUser, ok := user.(models.CurrentUser)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 记录历史申报
	baseSubmission := models.SubmitHistory{
		AuditorId: currentUser.UserId,
		Status:    information.Status,
		Cause:     information.Cause,
		Advice:    information.Advice,
	}
	var submissions []interface{}
	var successCount, failureCount int
	var errorMessage string
	for _, submissionId := range information.SubmissionIds {
		doc := baseSubmission
		doc.SubmissionId = submissionId
		submissions = append(submissions, doc)
	}

	insertResult, err := collection.InsertMany(context.TODO(), submissions)
	if err != nil {
		errorMessage = err.Error()
		log.Println("Insert error:", err)
	} else {
		// 计算成功和失败的个数
		successCount = len(insertResult.InsertedIDs)
		failureCount = len(information.SubmissionIds) - successCount
	}
	// 生成相应结构体
	response := auditManyResponse{
		SuccessCount: successCount,
		FailureCount: failureCount,
		Error:        errorMessage,
	}
	utils.ResponseSuccess(c, response)
}

type getAuditlist struct {
	Index          int `json:"index"`
	PaginationSize int `json:"paginationSize"`
}

func GetAuditHistory(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var information getAuditlist
	const DatabaseName string = ""
	const CollectionName string = "" //student

	err := c.ShouldBindJSON(&information)
	if err != nil {
		utils.ResponseError(c, "Paramter", "ParameterErrorMsg")
		return
	}
	// 获取collection
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	database := mongoClient.Database(DatabaseName)
	collection := database.Collection(CollectionName)
	filter := bson.D{}
	options := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip(int64(information.Index) * int64(information.PaginationSize)).SetLimit(int64(information.PaginationSize))
	result, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		// 处理逻辑
		return
	}
	var list []models.SubmitHistory
	if err = result.All(context.TODO(), &list); err != nil {
		// TODO: handle
		return
	}
	utils.ResponseSuccess(c, list)
	return
}
