package squarehandler

import (
	"hr/app/service"
	"hr/app/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type response struct {
	View int64 `json:"view"`
	Like int64 `json:"like"`
}

func GetViewsAndlikes(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	topicId := c.Param("topicID")
	views := service.GetTopicViews(c, topicId)
	likes := service.GetTopicLikes(c, topicId)

	utils.ResponseSuccess(c, &response{
		View: views,
		Like: likes,
	})
}

func LikesTopic(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	topicId := c.Param("topicID")
	filter := bson.M{
		"_id": topicId,
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
	replyId := c.Param("replyID")
	filter := bson.M{
		"_id": replyId,
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
