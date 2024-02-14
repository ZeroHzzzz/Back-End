package squarehandler

import (
	"context"
	"fmt"
	"hr/app/service"
	"hr/app/utils"
	"hr/configs/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 新评论
type newReplyInformation struct {
	ParentId string `json:"parentId"`
	Content  string `json:"content"`
}

func NewReply(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	topicId := c.Query("topicId")
	var information newReplyInformation
	err := c.ShouldBindJSON(&information)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}

	// 从上下文中获取用户信息
	currentUser := service.GetCurrentUser(c)
	// 新建submission记录
	newReply := models.Reply{
		ReplyId:   primitive.NewObjectID(),
		TopicId:   topicId,
		ParentId:  information.ParentId,
		Content:   information.Content,
		AutherId:  currentUser.UserId,
		CreatedAt: time.Now(),
	}
	_ = service.InsertOne(c, utils.MongodbName, utils.Reply, newReply)
	service.PublishMessage(c, utils.UserExchange, information.ParentId, utils.Reply2you) // 发布信息
	utils.ResponseSuccess(c, nil)
}

func GetReply(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	topicId := c.Query("topicId")
	pageParam := c.Query("page")
	page, err := strconv.Atoi(pageParam)
	limitParam := c.Query("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}
	fmt.Println(topicId)
	filter := bson.M{
		"topicId": topicId,
	}
	options := options.Find().SetSort(bson.D{{Key: "likes", Value: -1}, {Key: "createdAt", Value: -1}}).SetSkip(int64(page-1) * int64(limit)).SetLimit(int64(limit))
	var Reply []models.Reply
	cursor := service.Find(c, utils.MongodbName, utils.Reply, filter, options)
	if err := cursor.All(context.TODO(), &Reply); err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return
	}
	utils.ResponseSuccess(c, Reply)
}
