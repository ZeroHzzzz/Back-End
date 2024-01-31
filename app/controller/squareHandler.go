package controller

import (
	scoredatabase "hr/app/service/square"
	"hr/app/utils"
	"hr/configs/models/square"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateTopicInformation struct {
	UserId  string `json:"user_id" binding:"userId"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func CreateTopic(c *gin.Context) {
	var topicInformation CreateTopicInformation
	err := c.ShouldBindJSON(&topicInformation)
	if err != nil {
		utils.ResponseError(c, "Paramter", "ParameterErrorMsg")
		return
	}
	topicid, err := scoredatabase.CreateTopic(topicInformation.UserId, topicInformation.Title, topicInformation.Content)
	if err != nil {
		//处理逻辑
	}
	utils.ResponseSuccess(c, topicid)
	return
}

type GetTopicListInformation struct {
	Start int64 `json:"start" binding:"required"`
	End   int64 `json:"end" binding:"required"`
}

func GetTopicList(c *gin.Context) {
	const DatabaseName string = ""
	const CollectionName string = ""

	var gettopiclistinformation GetTopicListInformation
	err := c.ShouldBindJSON(&gettopiclistinformation)
	if err != nil {
		utils.ResponseError(c, "Paramter", "ParameterErrorMsg")
		return
	}

	// 获取collection
	var topiclist []square.Topic
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	database := mongoClient.Database(DatabaseName)
	collection := database.Collection(CollectionName)

	topiclist, err = scoredatabase.GetTopicList(gettopiclistinformation.Start, gettopiclistinformation.End, collection)
	if err != nil {
		//处理逻辑
	}
	utils.ResponseSuccess(c, topiclist)
	return
}
