package squarehandler

import (
	"fmt"
	"hr/app/service"
	"hr/app/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type response struct {
	View int64 `json:"View"`
	Like int64 `json:"Like"`
}

func GetViewsAndlikes(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	topicID := c.Query("TopicID")
	fmt.Println(topicID)
	views, likes := service.GetTopicViewsALikes(c, topicID)

	utils.ResponseSuccess(c, &response{
		View: views,
		Like: likes,
	})
}

func LikesTopic(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	topicID := c.Query("TopicID")
	objectID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return
	}
	filter := bson.M{
		"_id": objectID,
	}
	modified := bson.M{
		"$inc": bson.M{
			"Likes": 1,
		},
	}
	// 增加
	_ = service.UpdateOne(c, utils.MongodbName, utils.Topic, filter, modified)
	utils.ResponseSuccess(c, nil)
}

func LikeReply(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	replyID := c.Query("ReplyID")
	objectID, err := primitive.ObjectIDFromHex(replyID)
	if err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return
	}
	filter := bson.M{
		"_id": objectID,
	}
	modified := bson.M{
		"$inc": bson.M{
			"Likes": 1,
		},
	}
	// 增加
	_ = service.UpdateOne(c, utils.MongodbName, utils.Reply, filter, modified)
	utils.ResponseSuccess(c, nil)
}
