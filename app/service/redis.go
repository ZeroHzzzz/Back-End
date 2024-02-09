package service

import (
	"context"
	"fmt"
	"hr/app/utils"
	"hr/configs/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
)

func GetTopicViews(c *gin.Context, topicId string) int64 {
	// 获取文章浏览量，先从缓存找，然后找不到再去数据库找
	redisClient := GetRedisClint(c)
	view, err := redisClient.Get(context.Background(), fmt.Sprintf("%s_Topic_Views", topicId)).Int64()
	if err == redis.Nil {
		// 数据库找
		filter := bson.M{"_id": topicId}
		var topic models.Topic
		e := FindOne(c, utils.MongodbName, utils.Topic, filter).Decode(&topic)
		if e != nil {
			c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
			return -1
		}
		return int64(topic.Views)
	} else if err != nil {
		return -1
	}
	return view
}

func GetTopicLikes(c *gin.Context, topicId string) int64 {
	redisClient := GetRedisClint(c)
	like, err := redisClient.Get(context.Background(), fmt.Sprintf("%s_Topic_Likes", topicId)).Int64()
	if err != nil {
		return -1
	}
	return like
}
