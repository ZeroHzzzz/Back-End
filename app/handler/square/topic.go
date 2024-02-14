package squarehandler

import (
	"context"
	"fmt"
	"hr/app/service"
	"hr/app/utils"
	"hr/configs/models"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 创建文章
type CreateTopicInformation struct {
	UserId  string `json:"userId"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func NewTopic(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var topicInformation CreateTopicInformation
	err := c.ShouldBindJSON(&topicInformation)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}
	user := service.GetCurrentUser(c)
	if user.UserId != topicInformation.UserId {
		c.Error(utils.UNAUTHORIZED)
		c.Abort()
		return
	}
	newTopic := models.Topic{
		TopicID:   primitive.NewObjectID(),
		Title:     topicInformation.Title,
		Content:   topicInformation.Content,
		AutherID:  topicInformation.UserId,
		CreatedAt: time.Now(),
	}
	insertResult := service.InsertOne(c, utils.MongodbName, utils.Topic, newTopic)
	if insertResult == nil {
		log.Println("insertResult is nil")
		return
	}
	utils.ResponseSuccess(c, insertResult.InsertedID)
	return
}

// 文章列表
func GetTopicList(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	pageParam := c.Query("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}
	limitParam := c.Query("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}
	filter := bson.D{}
	options := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip((int64(page) - 1) * int64(limit)).SetLimit(int64(limit))
	result := service.Find(c, utils.MongodbName, utils.Topic, filter, options)
	var list []models.Topic
	if err = result.All(context.TODO(), &list); err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return
	}
	utils.ResponseSuccess(c, list)
	return
}

// 文章内容
func GetTopic(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	topicId := c.Query("topicId")
	fmt.Println(topicId)
	objectId, err := primitive.ObjectIDFromHex(topicId)
	if err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return
	}
	// 获取collection
	filter := bson.M{
		"_id": objectId,
	}
	var topic models.Topic
	result := service.FindOne(c, utils.MongodbName, utils.Topic, filter)
	if result == nil {
		c.Abort()
		return
	}
	err = result.Decode(&topic)
	if err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return
	}
	// 更新浏览量
	filter = bson.M{
		"_id": objectId,
	}
	modified := bson.M{
		"$inc": bson.M{
			"views": 1,
		},
	}
	_ = service.UpdateOne(c, utils.MongodbName, utils.Topic, filter, modified)
	utils.ResponseSuccess(c, topic)
	return
}

type ModifiedTopicInformation struct {
	Title   string `json:"title"`
	Context string `json:"context"`
}

func ModifiedTopic(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var information ModifiedTopicInformation
	err := c.ShouldBindJSON(&information)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}
	topicId := c.Query("topicId")
	objectId, err := primitive.ObjectIDFromHex(topicId)
	if err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return
	}
	// 从上下文中获取currentUser
	currentUser := service.GetCurrentUser(c)

	// 如果通过文章的id和修改人的id进行查找，如果找不到，说明修改人不是原作者，不允许修改
	filter := bson.M{
		"_id":      objectId,
		"autherId": currentUser.UserId,
	}
	modified := bson.M{
		"$set": bson.M{
			"title":   information.Title,
			"content": information.Context,
		},
	}
	service.UpdateOne(c, utils.MongodbName, utils.Topic, filter, modified)

	utils.ResponseSuccess(c, nil)
}
