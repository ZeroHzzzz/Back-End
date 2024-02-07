package service

import (
	"hr/app/utils"
	"hr/configs/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetCurrentUser(c *gin.Context) models.CurrentUser {
	user, ok := c.Get("currentUser")
	currentUser, ok := user.(models.CurrentUser)
	if !ok {
		c.Error(utils.GetError(utils.CONTEXT_ERROR, nil))
		c.Abort()
	}
	return currentUser
}

func GetmongoClient(c *gin.Context) *mongo.Client {
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.Error(utils.GetError(utils.CONTEXT_ERROR, nil))
		c.Abort()
		return nil
	}
	return mongoClient
}

func GetRabbitMQMiddle(c *gin.Context) *models.RabbitMQMiddleware {
	rabbitmqmiddle, exists := c.Request.Context().Value("RabbitMQMiddleware").(*models.RabbitMQMiddleware)
	if !exists {
		c.Error(utils.GetError(utils.CONTEXT_ERROR, nil))
		c.Abort()
		return nil
	}
	return rabbitmqmiddle
}

func GetRedisClint(c *gin.Context) *redis.Client {
	redisClient, exists := c.Request.Context().Value("redisClient").(*redis.Client)
	if !exists {
		c.Error(utils.GetError(utils.CONTEXT_ERROR, nil))
		c.Abort()
		return nil
	}
	return redisClient
}

// func Sync(c *gin.Context) {

// }
