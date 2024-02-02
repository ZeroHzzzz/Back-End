package square

import (
	"context"
	"fmt"
	"hr/app/utils"
	"hr/configs/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 创建文章
type CreateTopicInformation struct {
	UserId  string `json:"user_id" binding:"userId"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func NewTopic(c *gin.Context) {
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

// 文章列表
func GetTopicList(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	const DatabaseName string = ""
	const CollectionName string = ""
	lastViewParam := c.Query("lastview")
	lastView, err := strconv.Atoi(lastViewParam)
	if err != nil {
		// TODO
		return
	}
	limitParam := c.Query("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		// TODO
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
	options := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip(int64(lastView)).SetLimit(int64(limit))
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

// 文章内容
func GetTopic(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	topicId := c.Param("topicId")

	// 获取collection
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	const DatabaseName string = ""
	const CollectionName string = ""
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

type ModifiedTopicInformation struct {
	Context string `json:"context"`
}

func ModifiedTopic(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var information ModifiedTopicInformation
	err := c.ShouldBindJSON(&information)
	if err != nil {
		// TODO
		return
	}
	topicId := c.Param("topicId")
	// 从上下文中获取currentUser
	user, ok := c.Get("currentUser")
	currentUser, ok := user.(models.CurrentUser)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	const DatabaseName string = ""
	const CollectionName string = ""
	database := mongoClient.Database(DatabaseName)
	collection := database.Collection(CollectionName)
	// 如果通过文章的id和修改人的id进行查找，如果找不到，说明修改人不是原作者，不允许修改
	filter := bson.M{
		"topicId": topicId,
		"userId":  currentUser.UserId,
	}
	modified := bson.M{
		"$set": bson.M{
			"content": information.Context,
		},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, modified)
	if err != nil {
		//TODO
		return
	}
	utils.ResponseSuccess(c, nil)
}
