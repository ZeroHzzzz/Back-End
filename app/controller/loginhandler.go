package controller

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

type CurrentUser struct {
	UserId     string
	userName   string
	grade      int
	profession string
}

func LoginHandler(c *gin.Context) {
	const DatabaseName string = ""
	const CollectionName string = ""

	c.Header("Content-Type", "application/json")
	userid := c.PostForm("userId")
	password := c.PostForm("passWord")
	var user models.Student
	// err := database.DB.Where("userId=?", "passWord=?", userid, password).First(&user).Error

	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	database := mongoClient.Database(DatabaseName)
	collection := database.Collection(CollectionName)

	filter := bson.M{
		"userId":   userid,
		"passWord": password,
	}

	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		//处理逻辑
		if err == mongo.ErrNoDocuments {
			fmt.Println("No matching document found")
			return
		}
		log.Fatal(err)
		return
	}

	currentUser := CurrentUser{
		UserId:     userid,
		userName:   user.UserName,
		grade:      user.Grade,
		profession: user.Profession,
	}
	c.Set("CurrentUser", currentUser) //将用户信息储存到上下文
	utils.ResponseSuccess(c, currentUser)
	return
}
