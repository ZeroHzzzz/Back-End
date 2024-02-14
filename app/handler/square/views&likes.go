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
	View int64 `json:"view"`
	Like int64 `json:"like"`
}

func GetViewsAndlikes(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	topicId := c.Query("topicId")
	fmt.Println(topicId)
	views, likes := service.GetTopicViewsALikes(c, topicId)

	utils.ResponseSuccess(c, &response{
		View: views,
		Like: likes,
	})
}

func LikesTopic(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	topicId := c.Query("topicId")
	objectId, err := primitive.ObjectIDFromHex(topicId)
	if err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return
	}
	filter := bson.M{
		"_id": objectId,
	}
	modified := bson.M{
		"$inc": bson.M{
			"likes": 1,
		},
	}
	// 增加
	_ = service.UpdateOne(c, utils.MongodbName, utils.Topic, filter, modified)
	utils.ResponseSuccess(c, nil)
}

func LikeReply(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	replyId := c.Query("replyId")
	objectId, err := primitive.ObjectIDFromHex(replyId)
	if err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return
	}
	filter := bson.M{
		"_id": objectId,
	}
	modified := bson.M{
		"$inc": bson.M{
			"likes": 1,
		},
	}
	// 增加
	_ = service.UpdateOne(c, utils.MongodbName, utils.Reply, filter, modified)
	utils.ResponseSuccess(c, nil)
}
