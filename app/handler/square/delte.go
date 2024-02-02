package square

import (
	"context"
	"hr/app/utils"
	"hr/configs/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteTopic(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	topicId := c.Param("topicId")
	// 从上下文中获取用户信息
	currentUser, ok := c.Get("CurrentUser")
	if !ok {
		// TODO:
		return
	}
	user, ok := currentUser.(models.CurrentUser)
	if !ok {
		// TODO:
		return
	}
	// 从上下文中获取mongo客户端
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	// 要改
	database := mongoClient.Database("Form")
	collection := database.Collection("Submission")
	var err error
	// 辅导员拥有删除文章的能力
	if user.Role == "counsellor" {
		filter := bson.M{
			"_id": topicId,
		}
		_, err = collection.DeleteOne(context.TODO(), filter)
		// 删除评论

	} else if user.Role == "student" {
		filter := bson.M{
			"_id":      topicId,
			"autherID": user.UserId,
		}
		_, err = collection.DeleteOne(context.TODO(), filter)

	}
	if err != nil {
		// TODO:
		return
	}
	// 要改
	database = mongoClient.Database("Form")
	collection = database.Collection("Submission")
	filter := bson.M{
		"topicId": topicId,
	}
	_, err = collection.DeleteMany(context.TODO(), filter)
	// 这里要注意有可能这个文章没有评论
	if err != nil {
		// TODO
		return
	}
	utils.ResponseSuccess(c, nil)
}
