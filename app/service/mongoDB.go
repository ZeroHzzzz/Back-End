package service

import (
	"context"
	"fmt"
	"hr/app/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoDB
// Update
func UpdateOne(c *gin.Context, databaseName, collectionName string, filter, modified interface{}) *mongo.UpdateResult {
	mongoClient := GetmongoClient(c)
	database := mongoClient.Database(databaseName)
	collection := database.Collection(collectionName)
	updateResult, err := collection.UpdateOne(context.Background(), filter, modified)
	// 如果没有合适的修改也会抛出错误
	if err != nil || updateResult.ModifiedCount == 0 {
		c.Error(utils.GetError(utils.DATABASE_OPERATION_ERROR, "UpdateOne Failed"))
		c.Abort()
		return nil
	}
	return updateResult
}

func UpdateMany(c *gin.Context, databaseName, collectionName string, filter, modified interface{}) *mongo.UpdateResult {
	mongoClient := GetmongoClient(c)
	database := mongoClient.Database(databaseName)
	collection := database.Collection(collectionName)
	updateResult, err := collection.UpdateMany(context.Background(), filter, modified)
	if err != nil {
		c.Error(utils.GetError(utils.DATABASE_OPERATION_ERROR, err.Error()))
		c.Abort()
		return nil
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
		c.Error(utils.GetError(utils.DATABASE_OPERATION_ERROR, err.Error()))
		c.Abort()
		return nil
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
		msg := fmt.Sprintf("Inserted %d , failed %d", successCount, failureCount)

		c.Error(utils.GetError(utils.DATABASE_OPERATION_ERROR, msg))
		c.Abort()
		return nil
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
		c.Error(utils.GetError(utils.DATABASE_OPERATION_ERROR, cursor.Err().Error()))
		c.Abort()
		return nil
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
		c.Error(utils.GetError(utils.DATABASE_OPERATION_ERROR, err.Error()))
		c.Abort()
		return nil
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
		c.Error(utils.GetError(utils.DATABASE_OPERATION_ERROR, err.Error()))
		c.Abort()
		return nil
	}
	return cursor
}

func DeleteMany(c *gin.Context, databaseName, collectionName string, filter interface{}) *mongo.DeleteResult {
	mongoClient := GetmongoClient(c)
	database := mongoClient.Database(databaseName)
	collection := database.Collection(collectionName)
	cursor, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		c.Error(utils.GetError(utils.DATABASE_OPERATION_ERROR, err.Error()))
		c.Abort()
		return nil
	}
	return cursor
}
