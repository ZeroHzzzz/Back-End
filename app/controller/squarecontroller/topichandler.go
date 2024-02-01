package squarecontroller

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

type CreateTopicInformation struct {
	UserId  string `json:"user_id" binding:"userId"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func CreateTopic(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	const DatabaseName string = ""
	const CollectionName string = "" //student

	var topicInformation CreateTopicInformation
	err := c.ShouldBindJSON(&topicInformation)
	if err != nil {
		utils.ResponseError(c, "Paramter", "ParameterErrorMsg")
		return
	}
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	database := mongoClient.Database(DatabaseName)
	collection := database.Collection(CollectionName)
	newTopic := models.Topic{
		Title:    topicInformation.Title,
		Content:  topicInformation.Content,
		AutherId: topicInformation.UserId,
	}
	insertResult, err := collection.InsertOne(context.TODO(), newTopic)
	if err != nil {
		log.Println("Insert error:", err)
		return
	}
	utils.ResponseSuccess(c, insertResult.InsertedID)
	return
}

type GetTopicListInformation struct {
	Index          int64 `json:"index" binding:"required"`
	PaginationSize int64 `json:"paginationSize"`
}

func GetTopicList(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	const DatabaseName string = ""
	const CollectionName string = ""

	var gettopiclistinformation GetTopicListInformation
	err := c.ShouldBindJSON(&gettopiclistinformation)
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
	options := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip(gettopiclistinformation.Index * gettopiclistinformation.PaginationSize).SetLimit(gettopiclistinformation.PaginationSize)
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

func GetTopic(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	const DatabaseName string = ""
	const CollectionName string = ""
	topicId := c.Param("topicId")

	// 获取collection
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	database := mongoClient.Database(DatabaseName)
	collection := database.Collection(CollectionName)

	filter := bson.M{
		"_id": topicId,
	}
	var topic models.Topic
	err := collection.FindOne(context.TODO(), filter).Decode(&topic)
	if err != nil {
		//处理逻辑
		if err == mongo.ErrNoDocuments {
			fmt.Println("No matching document found")
			return
		}
		log.Fatal(err)
		return
	}
	utils.ResponseSuccess(c, topic)
	return
}
