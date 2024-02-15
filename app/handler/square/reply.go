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
	ParentID string `json:"ParentID"`
	Content  string `json:"Content"`
}

func NewReply(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	topicID := c.Query("TopicID")
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
		ReplyID:  primitive.NewObjectID(),
		TopicID:  topicID,
		ParentID: information.ParentID,
		Content:  information.Content,
		AutherID: currentUser.UserID,
		CreateAt: time.Now().Unix(),
	}
	_ = service.InsertOne(c, utils.MongodbName, utils.Reply, newReply)
	service.PublishMessage(c, utils.UserExchange, information.ParentID, fmt.Sprintf("%s %s", currentUser.UserID, utils.Reply2you)) // 发布信息
	utils.ResponseSuccess(c, nil)
}

func GetReply(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	topicID := c.Query("TopicID")
	pageParam := c.Query("Page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}
	limitParam := c.Query("Limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}
	fmt.Println(topicID)
	filter := bson.M{
		"TopicID": topicID,
	}
	options := options.Find().SetSort(bson.D{{Key: "Likes", Value: -1}, {Key: "CreateAt", Value: -1}}).SetSkip(int64(page-1) * int64(limit)).SetLimit(int64(limit))
	var Reply []models.Reply
	cursor := service.Find(c, utils.MongodbName, utils.Reply, filter, options)
	if err := cursor.All(context.TODO(), &Reply); err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return
	}
	utils.ResponseSuccess(c, Reply)
}
