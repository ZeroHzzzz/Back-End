package service

import (
	"hr/app/utils"
	"hr/configs/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTopicViewsALikes(c *gin.Context, topicID string) (int64, int64) {
	objectID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return -1, -1
	}
	// 获取文章浏览量，先从缓存找，然后找不到再去数据库找
	filter := bson.M{"_id": objectID}
	var topic models.Topic
	result := FindOne(c, utils.MongodbName, utils.Topic, filter)
	if result.Err() != nil {
		c.Error(utils.GetError(utils.MONGODB_OPERATION_ERROR, result.Err().Error()))
		return -1, -1
	}
	decodeErr := result.Decode(&topic)
	if decodeErr != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, decodeErr.Error()))
		return -1, -1
	}
	return int64(topic.Views), int64(topic.Likes)
}
