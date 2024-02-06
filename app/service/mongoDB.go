package service

import (
	"context"
	"fmt"
	"hr/app/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// const (
// 	MongoDBHost     = "localhost"
// 	MongoDBPort     = 27017
// 	MongoDBPassword = "password" // 这里没有配置密码，很有可能出问题
// )

// // 初始化 MongoDB 客户端
// func InitMongoClient(c *gin.Context) *mongo.Client {
// 	// 设置 MongoDB 连接配置
// 	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", MongoDBHost, MongoDBPort)).SetConnectTimeout(10 * time.Second)

// 	// 连接 MongoDB
// 	client, err := mongo.Connect(context.Background(), clientOptions)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize MongoDB client"})
// 		c.Abort()
// 		return nil
// 	}

// 	// 设置最大连接池大小
// 	clientOptions.SetMaxPoolSize(10)

// 	// 创建 MongoDB 连接上下文
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel() // 这里会出问题的，应该放在handler里面而不是放在函数内部

// 	// 检查连接
// 	err = client.Ping(ctx, nil)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize MongoDB client"})
// 		c.Abort()
// 		return nil
// 	}
// 	fmt.Println("Connected to MongoDB!")
// 	return client
// }

// func CloseMongoClient(client *mongo.Client) {
// 	if client != nil {
// 		err := client.Disconnect(context.Background())
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println("MongoDB connection closed.")
// 	}
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
