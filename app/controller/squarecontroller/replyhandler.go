package squarecontroller

import (
	"context"
	"hr/app/utils"
	"hr/configs/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type replyTopic struct {
	ParentId string `json:"parentId"`
	Content  string `json:"content"`
}

func ReplyTopic(c *gin.Context) {
	const DatabaseName string = ""
	const CollectionName string = ""
	c.Header("Content-Type", "application/json")
	topicId := c.Param("topicId")
	var information replyTopic
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
		ParertId: information.ParentId,
		Content:  information.Content,
		UserId:   user.UserId,
	}
	_, err = collection.InsertOne(context.TODO(), newReply)
	if err != nil {
		// TODO:
		return
	}
	utils.ResponseSuccess(c, nil)
}
