package service

import (
	"context"
	"fmt"
	"hr/app/utils"
	"hr/configs/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCurrentUser(c *gin.Context) models.CurrentUser {
	user, ok := c.Get("currentUser")
	currentUser, ok := user.(models.CurrentUser)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	return currentUser
}

func GetmongoClient(c *gin.Context) *mongo.Client {
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.Error(utils.GetError(utils.DATABASE_ERROR, "mongoClient init failed"))
		c.Abort()
	}
	return mongoClient
}

func GetRabbitMQMiddle(c *gin.Context) *models.RabbitMQMiddleware {
	rabbitmqmiddle, exists := c.Request.Context().Value("RabbitMQMiddleware").(*models.RabbitMQMiddleware)
	if !exists {
		c.Error(utils.GetError(utils.DATABASE_ERROR, "RabbitMQMiddleware init failed"))
		c.Abort()
	}
	return rabbitmqmiddle
}

// func GetRedisClint(c *gin.Context) *redis.Client {

// }
// mongoDB
// Update
func UpdateOne(c *gin.Context, databaseName, collectionName string, filter, modified interface{}) *mongo.UpdateResult {
	mongoClient := GetmongoClient(c)
	database := mongoClient.Database(databaseName)
	collection := database.Collection(collectionName)
	updateResult, err := collection.UpdateOne(context.Background(), filter, modified)
	// 如果没有合适的修改也会抛出错误
	if err != nil || updateResult.ModifiedCount == 0 {
		c.Error(utils.GetError(utils.DATABASE_ERROR, "Operation Failed"))
		c.Abort()
	}
	return updateResult
}

func UpdateMany(c *gin.Context, databaseName, collectionName string, filter, modified interface{}) *mongo.UpdateResult {
	mongoClient := GetmongoClient(c)
	database := mongoClient.Database(databaseName)
	collection := database.Collection(collectionName)
	updateResult, err := collection.UpdateMany(context.Background(), filter, modified)
	if err != nil {
		c.Error(utils.GetError(utils.DATABASE_ERROR, "Operation Failed"))
		c.Abort()
	}
	return updateResult
}

// Insert
func InsertOne(c *gin.Context, databaseName, collectionName string, document interface{}) *mongo.InsertOneResult {
	mongoClient := GetmongoClient(c)
	database := mongoClient.Database(databaseName)
	collection := database.Collection(collectionName)
	insertOneResult, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		c.Error(utils.GetError(utils.DATABASE_ERROR, "Operation Failed"))
		c.Abort()
	}
	return insertOneResult
}

func InsertMany(c *gin.Context, databaseName, collectionName string, document []interface{}) *mongo.InsertManyResult {
	mongoClient := GetmongoClient(c)
	database := mongoClient.Database(databaseName)
	collection := database.Collection(collectionName)
	insertManyResult, err := collection.InsertMany(context.Background(), document)
	if err != nil {
		successCount := len(insertManyResult.InsertedIDs)
		failureCount := len(document) - successCount
		tmp := fmt.Errorf("Inserted %d , failed %d", successCount, failureCount)
		myErr := utils.DATABASE_ERROR
		myErr.Data = tmp.Error()
		c.Error(myErr)
		c.Abort()
	}
	return insertManyResult
}

// Find
func FindOne(c *gin.Context, databaseName, collectionName string, filter interface{}) *mongo.SingleResult {
	mongoClient := GetmongoClient(c)
	database := mongoClient.Database(databaseName)
	collection := database.Collection(collectionName)
	cursor := collection.FindOne(c, filter)
	if cursor.Err() != nil {
		c.Error(utils.GetError(utils.DATABASE_ERROR, cursor.Err().Error()))
		c.Abort()
	}
	return cursor
}

func Find(c *gin.Context, databaseName, collectionName string, filter interface{}, opts ...*options.FindOptions) *mongo.Cursor {
	mongoClient := GetmongoClient(c)
	database := mongoClient.Database(databaseName)
	collection := database.Collection(collectionName)
	fo := options.MergeFindOptions(opts...)
	cursor, err := collection.Find(context.Background(), filter, fo)
	if err != nil {
		c.Error(utils.GetError(utils.DATABASE_ERROR, "Operation Failed"))
		c.Abort()
	}
	return cursor
}

// Delete
func DeleteOne(c *gin.Context, databaseName, collectionName string, filter interface{}) *mongo.DeleteResult {
	mongoClient := GetmongoClient(c)
	database := mongoClient.Database(databaseName)
	collection := database.Collection(collectionName)
	cursor, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		c.Error(utils.GetError(utils.DATABASE_ERROR, err.Error()))
		c.Abort()
	}
	return cursor
}

func DeleteMany(c *gin.Context, databaseName, collectionName string, filter interface{}) *mongo.DeleteResult {
	mongoClient := GetmongoClient(c)
	database := mongoClient.Database(databaseName)
	collection := database.Collection(collectionName)
	cursor, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		c.Error(utils.GetError(utils.DATABASE_ERROR, err.Error()))
		c.Abort()
	}
	return cursor
}

// rmq
func DeclareQueue(c *gin.Context, queueName string) amqp.Queue {
	r := GetRabbitMQMiddle(c)
	q, err := r.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		true,      // no-wait
		nil,       // arguments
	)
	if err != nil {
		c.Error(err)
		c.Abort()
	}
	return q
}

func DeclareExchange(c *gin.Context, exchangeName string) {
	// 声明交换机
	r := GetRabbitMQMiddle(c)
	err := r.Channel.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.Error(err)
		c.Abort()
	}
}

func PublishMessage(c *gin.Context, exchangeName, queueName, message string) {
	r := GetRabbitMQMiddle(c)
	err := r.Channel.Publish(
		exchangeName, // exchange
		queueName,    // routing key (队列名即为路由键) ,如果为空就是发布到全部队列
		false,        // mandatory
		false,        //
		//  immediate 参数为 false 时（默认值），如果消息无法被立即路由到队列，消息将会被存储在队列中等待消费者接收
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(message),
			DeliveryMode: amqp.Persistent, // 持久化
		},
	)
	if err != nil {
		c.Error(err)
		c.Abort()
	}
}

func ConsumeMessage(c *gin.Context, queueName string) <-chan amqp.Delivery {
	// 声明一个消费者，并返回一个接收（receive）操作符用于从通道中接收数据的表达式
	r := GetRabbitMQMiddle(c)
	msgs, err := r.Channel.Consume(
		queueName, // queueName
		"",        // consumer
		false,     //autoAck 自动确认已读，false为手动
		false,     // exclusive 独占
		true,      // 接收自己的信息
		true,      // 不等待服务器响应，false表示等待
		nil,       // 其他参数
	)
	if err != nil {
		c.Error(err)
		c.Abort()
	}
	return msgs
}

// for msg := range messages {
// 	// 将消息发送到前端
// 	sendMessageToClient(msg.Body)
// }
