package squarecontroller

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

type response struct {
	topic models.Topic
	reply []models.Reply
}

func GetTopic(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	topicId := c.Param("topicId")

	// 获取collection
	// 这里还要返回这条topic的回复，按照点赞量排序
	// 其实如果用分页加载的话，只用找点赞前几名的回复就行了
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	// 找文章
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
	// 找回复
	// const DatabaseName string = ""
	// const CollectionName string = ""
	var reply []models.Reply
	database = mongoClient.Database(DatabaseName)
	collection = database.Collection(CollectionName)
	filter = bson.M{
		"_id": topicId,
	}
	options := options.Find().SetSort(bson.D{{Key: "likes", Value: -1}}).SetLimit(25) //第一次显示25条评论
	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		// 处理
		return
	}
	if err := cursor.All(context.TODO(), &reply); err != nil {
		return
	}
	response := response{
		topic: topic,
		reply: reply,
	}
	utils.ResponseSuccess(c, response)
	return
}

// 这是用来分页加载的
func GetFatherReply(c *gin.Context) {
	// 加载父评论
	c.Header("Content-Type", "application/json")
	topicId := c.Param("topicId")
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
	// 获取客户端
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
		"_id":     topicId,
		"partent": "",
	}
	options := options.Find().SetSort(bson.D{{Key: "likes", Value: -1}}).SetSkip(int64(lastView)).SetLimit(int64(limit))
	var parentReply []models.Reply
	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		// 处理
		return
	}
	if err := cursor.All(context.TODO(), &parentReply); err != nil {
		return
	}
	utils.ResponseSuccess(c, parentReply)
}
func GetSonReply(c *gin.Context) {
	// /square/{topicId}/reply
	// 加载子评论
	c.Header("Content-Type", "application/json")

	topicId := c.Param("topicId")
	replyIdParams := c.Query("replyId")
	replyId, err := strconv.Atoi(replyIdParams)
	if err != nil {
		// TODO
		return
	}
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
	// 获取客户端
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
		"_id":    topicId,
		"parent": replyId,
	}
	options := options.Find().SetSort(bson.D{{Key: "likes", Value: -1}}).SetSkip(int64(lastView)).SetLimit(int64(limit))
	var parentReply []models.Reply
	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		// 处理
		return
	}
	if err := cursor.All(context.TODO(), &parentReply); err != nil {
		return
	}
	utils.ResponseSuccess(c, parentReply)
}
