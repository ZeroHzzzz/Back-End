package service

import (
	"hr/app/utils"
	"hr/configs/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetCurrentUser(c *gin.Context) models.CurrentUser {
	user, ok := c.Get("CurrentUser")
	currentUser, ok := user.(models.CurrentUser)
	if !ok {
		c.Error(utils.UNAUTHORIZED)
		c.Abort()
	}
	return currentUser
}

func GetmongoClient(c *gin.Context) *mongo.Client {
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.Error(utils.MONGODB_INIT_ERROR)
		c.Abort()
		return nil
	}
	return mongoClient
}

func GetRabbitMQMiddle(c *gin.Context) *models.RabbitMQMiddleware {
	rabbitmqMiddleware, exists := c.Get("RabbitMQMiddleware")
	if !exists {
		log.Println("没有找到 RabbitMQMiddleware 上下文")
		return nil
	}
	rmqMiddleware, ok := rabbitmqMiddleware.(*models.RabbitMQMiddleware)
	if !ok {
		log.Println("无法转换为 RabbitMQMiddleware")
		return nil
	}
	return rmqMiddleware
}

func GetRedisClint(c *gin.Context) *redis.Client {
	redisClient, exists := c.Request.Context().Value("redisClient").(*redis.Client)
	if !exists {
		c.Error(utils.REDIS_INIT_ERROR)
		c.Abort()
		return nil
	}
	return redisClient
}

// func Sync(c *gin.Context) {

// }
