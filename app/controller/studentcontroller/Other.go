package studentcontroller

import (
	"context"
	"hr/app/utils"
	"hr/configs/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type information struct {
	UserId   string `json:"userId"`
	Category string `json:"category"`
	Content  string `json:"content"`
}

func Feedback(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var information information
	const DatabaseName string = ""
	const CollectionName string = ""
	err := c.ShouldBindJSON(&information)
	if err != nil {
		utils.ResponseError(c, "Paramter", "ParameterErrorMsg")
		return
	}
	// 从上下文中获取mongo客户端
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}

	// feedback
	database := mongoClient.Database(DatabaseName)
	collection := database.Collection(CollectionName)

	newFeeback := models.Feedback{
		Category: information.Category,
		UserId:   information.UserId,
		Content:  information.Content,
		Status:   false,
	}
	insertResult, err := collection.InsertOne(context.Background(), newFeeback)
	if err != nil {
		log.Fatal(err)
	}

	// student更新feedback列表
	// database = mongoClient.Database(DatabaseName)
	// collection = database.Collection("")
	// filter := bson.M{
	// 	"userId":   information.UserId,
	// 	"category": information.Category,
	// }
	// modified := bson.M{
	// 	"$push": bson.M{
	// 		"FeedbackId": insertResult.InsertedID,
	// 	},
	// }
	// _, err = collection.UpdateOne(context.TODO(), filter, modified)
	// if err != nil {
	// 	//处理逻辑
	// 	if err == mongo.ErrNoDocuments {
	// 		fmt.Println("No matching document found")
	// 		return
	// 	}
	// 	log.Fatal(err)
	// 	return
	// }
	utils.ResponseSuccess(c, insertResult.InsertedID) //返回文档的id
	return
}

func Advice(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var information information
	const DatabaseName string = ""
	const CollectionName string = ""
	err := c.ShouldBindJSON(&information)
	if err != nil {
		utils.ResponseError(c, "Paramter", "ParameterErrorMsg")
		return
	}
	// 从上下文中获取mongo客户端
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}

	// feedback
	database := mongoClient.Database(DatabaseName)
	collection := database.Collection(CollectionName)

	newFeeback := models.Feedback{
		Category: information.Category,
		UserId:   information.UserId,
		Content:  information.Content,
		Status:   false,
	}
	insertResult, err := collection.InsertOne(context.Background(), newFeeback)
	if err != nil {
		log.Fatal(err)
	}

	// student更新feedback列表
	// database = mongoClient.Database(DatabaseName)
	// collection = database.Collection("")
	// filter := bson.M{
	// 	"userId":   information.UserId,
	// 	"category": information.Category,
	// }
	// modified := bson.M{
	// 	"$push": bson.M{
	// 		"AdviceId": insertResult.InsertedID,
	// 	},
	// }
	// _, err = collection.UpdateOne(context.TODO(), filter, modified)
	// if err != nil {
	// 	//处理逻辑
	// 	if err == mongo.ErrNoDocuments {
	// 		fmt.Println("No matching document found")
	// 		return
	// 	}
	// 	log.Fatal(err)
	// 	return
	// }
	utils.ResponseSuccess(c, insertResult.InsertedID) //返回文档的id
	return
}
