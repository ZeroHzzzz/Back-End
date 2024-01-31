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
)

type auditSubmissionInformation struct {
	Status bool   `json:"status"`
	Cause  string `json:"cause"`
	Advice string `json:"advice"`
}

func AuditSubmission(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	const DatabaseName string = ""
	const CollectionName string = "" //submission
	var information auditSubmissionInformation
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
