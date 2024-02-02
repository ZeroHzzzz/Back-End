package square

import (
	"context"
	"hr/app/utils"
	"hr/configs/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
		utils.ResponseError(c, "failure", "Parameter wrong")
		return
	}

	// 从上下文中获取用户信息
	currentUser, ok := c.Get("CurrentUser")
	if !ok {
		// TODO:
		return
	}
	user, ok := currentUser.(models.CurrentUser)
	if err != nil {
		// TODO:
		return
	}

	// 从上下文中获取mongo客户端
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	// 新建submission记录
	database := mongoClient.Database("Form")
	collection := database.Collection("Submission")

	newReply := models.Reply{
		TopicId:  topicId,
		ParentId: information.ParentId,
		Content:  information.Content,
		AutherId: user.UserId,
	}
	_, err = collection.InsertOne(context.TODO(), newReply)
	if err != nil {
		// TODO:
		return
	}
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
		"_id": topicId,
	}
	options := options.Find().SetSort(bson.D{{Key: "likes", Value: -1}}).SetSkip(int64(page) * int64(limit)).SetLimit(int64(limit))
	var Reply []models.Reply
	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		// 处理
		return
	}
	if err := cursor.All(context.TODO(), &Reply); err != nil {
		return
	}
	utils.ResponseSuccess(c, Reply)
}
