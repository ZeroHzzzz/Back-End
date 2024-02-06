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

// 新评论
type newReplyInformation struct {
	ParentId string `json:"parentId"`
	Content  string `json:"content"`
}

func NewReply(c *gin.Context) {
	const DatabaseName string = ""
	const CollectionName string = ""
	c.Header("Content-Type", "application/json")
	topicId := c.Param("topicId")
	var information newReplyInformation
	err := c.ShouldBindJSON(&information)
	if err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}

	// 从上下文中获取用户信息
	currentUser := service.GetCurrentUser(c)
	// 新建submission记录
	newReply := models.Reply{
		TopicId:  topicId,
		ParentId: information.ParentId,
		Content:  information.Content,
		AutherId: currentUser.UserId,
	}
	_ = service.InsertOne(c, "", "", newReply)
	service.PublishMessage(c, utils.UserExchange, information.ParentId, utils.Reply2you) // 发布信息
	utils.ResponseSuccess(c, nil)
}

func GetReply(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	topicId := c.Param("topicId")
	pageParam := c.Query("lastview")
	page, err := strconv.Atoi(pageParam)
	limitParam := c.Query("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}

	filter := bson.M{
		"_id": topicId,
	}
	options := options.Find().SetSort(bson.D{{Key: "likes", Value: -1}}).SetSkip(int64(page) * int64(limit)).SetLimit(int64(limit))
	var Reply []models.Reply
	cursor := service.Find(c, "", "", filter, options)
	if err := cursor.All(context.TODO(), &Reply); err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}
	utils.ResponseSuccess(c, Reply)
}
