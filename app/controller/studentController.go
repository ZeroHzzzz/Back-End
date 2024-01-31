package controller

import (
	"context"
	"hr/app/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type submitInformation struct {
	CurrentUser
	itemName     string
	academicYear string
	evidence     []string
}

const savePath = ""

func ScoreHandler(c *gin.Context) {
	// 上传申报
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
	if currentUser.(CurrentUser).UserId != userId {
		return
	}
	newSubmission := submitInformation{
		CurrentUser:  currentUser.(CurrentUser),
		itemName:     itemName,
		academicYear: academicYear,
		evidence:     destPaths,
	}

	// 从上下文中获取mongo客户端
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	database := mongoClient.Database("your_database_name")
	collection := database.Collection("your_collection_name")

	insertResult, err := collection.InsertOne(context.Background(), newSubmission)
	if err != nil {
		log.Fatal(err)
	}
	utils.ResponseSuccess(c, insertResult.InsertedID)
}
