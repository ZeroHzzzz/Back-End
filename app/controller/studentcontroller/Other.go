package studentcontroller

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

type feedbackInformation struct {
	UserId  string `json:"userId"`
	Content string `json:"content"`
}

func Feedback(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var feedbackinformation feedbackInformation
	const DatabaseName string = ""
	const CollectionName string = ""
	err := c.ShouldBindJSON(&feedbackinformation)
	if err != nil {
		utils.ResponseError(c, "Paramter", "ParameterErrorMsg")
		return
	}
	// 从上下文中获取mongo客户端
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}

	// feedback
	database := mongoClient.Database(DatabaseName)
	collection := database.Collection(CollectionName)

	newFeeback := models.Feedback{
		Category: "Feedback",
		UserId:   feedbackinformation.UserId,
		Content:  feedbackinformation.Content,
		Status:   false,
	}
	insertResult, err := collection.InsertOne(context.Background(), newFeeback)
	if err != nil {
		log.Fatal(err)
	}

	// student更新feedback列表
	database = mongoClient.Database(DatabaseName)
	collection = database.Collection("")
	filter := bson.M{
		"userId":   feedbackinformation.UserId,
		"category": "Feedback",
	}
	modified := bson.M{
		"$push": bson.M{
			"FeedbackId": insertResult.InsertedID,
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
	utils.ResponseSuccess(c, insertResult.InsertedID) //返回文档的id
	return
}
