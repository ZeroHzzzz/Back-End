package square

import (
	"context"
	"hr/app/service"
	"hr/app/utils"
	"hr/configs/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 创建文章
type CreateTopicInformation struct {
	UserId  string `json:"userID" binding:"userId"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func NewTopic(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	var topicInformation CreateTopicInformation
	err := c.ShouldBindJSON(&topicInformation)
	if err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}
	newTopic := models.Topic{
		Title:    topicInformation.Title,
		Content:  topicInformation.Content,
		AutherID: topicInformation.UserId,
	}
	insertResult := service.InsertOne(c, "", "", newTopic)
	utils.ResponseSuccess(c, insertResult.InsertedID)
	return
}

// 文章列表
func GetTopicList(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	pageParam := c.Query("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}
	limitParam := c.Query("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}
	filter := bson.D{}
	options := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip((int64(page) - 1) * int64(limit)).SetLimit(int64(limit))
	result := service.Find(c, "", "", filter, options)
	var list []models.SubmitHistory
	if err = result.All(context.TODO(), &list); err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
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
	filter := bson.M{
		"_id": topicId,
	}
	var topic models.Topic
	err := service.FindOne(c, "", "", filter).Decode(&topic)
	if err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
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
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}
	topicId := c.Param("topicId")

	// 从上下文中获取currentUser
	currentUser := service.GetCurrentUser(c)

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
	_ = service.UpdateOne(c, "", "", filter, modified)

	utils.ResponseSuccess(c, nil)
}
