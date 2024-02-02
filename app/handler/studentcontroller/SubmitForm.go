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

const savePath = ""

func SubmitHandler(c *gin.Context) {
	// 上传申报
	const DatabaseName string = ""
	const CollectionName string = ""
	c.Header("Content-Type", "application/json")
	userId := c.Param("userId")
	itemName := c.PostForm("itemName")
	academicYear := c.PostForm("academicYear")
	data, err := c.MultipartForm()
	if err != nil {
		return
	}
	files := data.File["evidence"]
	destPaths := make([]string, len(files))

	for i, file := range files {
		dst := savePath + "/" + file.Filename
		destPaths[i] = dst
		err := c.SaveUploadedFile(file, dst)
		if err != nil {
			return
		}
	}

	// 从上下文中获取用户信息
	currentUser, ok := c.Get("CurrentUser")
	if !ok {
		return
	}
	// 获取的currentUser需要断言
	if currentUser.(models.CurrentUser).UserId != userId {
		return
	}

	// 从上下文中获取mongo客户端
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	// 新建submission记录
	database := mongoClient.Database("Form")
	collection := database.Collection("Submission")
	newSubmission := models.SubmitInformation{
		CurrentUser:  currentUser.(models.CurrentUser),
		ItemName:     itemName,
		AcademicYear: academicYear,
		Evidence:     destPaths,
		Status:       false,
	}
	insertResult, err := collection.InsertOne(context.Background(), newSubmission)
	if err != nil {
		log.Fatal(err)
	}

	// student更新submission列表
	// database = mongoClient.Database(DatabaseName)
	// collection = database.Collection("")
	// filter := bson.M{
	// 	"userId": userId,
	// }
	// modified := bson.M{
	// 	"$push": bson.M{
	// 		"SubmissionId": insertResult.InsertedID,
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
	utils.ResponseSuccess(c, insertResult.InsertedID)
}

func GetSubmissionStatus(c *gin.Context) {
	// 从Form数据库中查找，然后返回每个的状态，数据格式应该是一个字典
	const DatabaseName string = ""
	const CollectionName string = ""
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

	result, err := collection.Find(context.TODO(), filter)
	if err != nil {
		//处理逻辑
		if err == mongo.ErrNoDocuments {
			fmt.Println("No matching document found")
			return
		}
		log.Fatal(err)
		return
	}
	var forms []models.SubmitInformation
	if err := result.All(context.Background(), &forms); err != nil {
		log.Fatal(err)
	}
	utils.ResponseSuccess(c, forms)
	return

}
