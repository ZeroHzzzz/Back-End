package studentcontroller

import (
	"context"
	"fmt"
	"hr/app/utils"
	"hr/configs/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetConcreteSorceHandler(c *gin.Context) {
	const DatabaseName string = ""
	const CollectionName string = ""
	// 上传申报
	c.Header("Content-Type", "application/json")
	userId := c.Param("userId")
	academicYear := c.Param("academicYear")

	// 从上下文中获取mongo客户端
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	database := mongoClient.Database(DatabaseName)
	collection := database.Collection(CollectionName)
	filter := bson.M{
		"userId":       userId,
		"academicTear": academicYear,
	}

	var result models.Score
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		//处理逻辑
		if err == mongo.ErrNoDocuments {
			fmt.Println("No matching document found")
			return
		}
		log.Fatal(err)
		return
	}
	utils.ResponseSuccess(c, result)
	return

}

func GetYearScoreHandler(c *gin.Context) {
	const DatabaseName string = ""
	const CollectionName string = ""
	// 上传申报
	c.Header("Content-Type", "application/json")
	userId := c.Param("userId")

	// 从上下文中获取mongo客户端
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	database := mongoClient.Database(DatabaseName)
	collection := database.Collection(CollectionName)
	filter := bson.M{
		"userId": userId,
	}

	var student models.Student
	err := collection.FindOne(context.TODO(), filter).Decode(&student)
	if err != nil {
		//处理逻辑
		if err == mongo.ErrNoDocuments {
			fmt.Println("No matching document found")
			return
		}
		log.Fatal(err)
		return
	}
	utils.ResponseSuccess(c, student.Grade)
	return

}
