package service

import (
	"context"
	"fmt"
	"hr/app/utils"
	"hr/configs/models"
	"net/http"

	"github.com/gin-gonic/gin"
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
