package studentcontroller

import (
	"context"
	"fmt"
	"hr/app/midware"
	"hr/app/utils"
	"hr/configs/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type modifiedprofileInformation struct {
	PassWord    string `json:"passWord"`
	NewPassword string `json:"newPassword"`
}

func ModifiedProfileHandler(c *gin.Context) {
	//修改密码后token行为可能需要深入考虑
	c.Header("Content-Type", "application/json")

	const DatabaseName string = ""
	const CollectionName string = "" //student

	var modifiedprofileinformation modifiedprofileInformation
	err := c.ShouldBindJSON(&modifiedprofileinformation)
	if err != nil {
		utils.ResponseError(c, "failure", "Parameter wrong")
		return
	}

	userid := c.Param("userId")

	// 从上下文中获取mongo客户端
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	database := mongoClient.Database("DatabaseName")
	collection := database.Collection("CollectionName")

	filter := bson.M{
		"userId":   userid,
		"passWord": modifiedprofileinformation.PassWord,
	}
	modified := bson.M{
		"$set": bson.M{
			"passWord": modifiedprofileinformation.NewPassword,
		},
	}
	// 修改之后的文档
	_, err = collection.UpdateOne(context.TODO(), filter, modified)
	if err != nil {
		//处理逻辑
		if err == mongo.ErrNoDocuments {
			fmt.Println("No matching document found")
			return
		}
		log.Fatal(err)
		return
	}
	// 从上下文中获取currentUser
	user, ok := c.Get("currentUser")
	currentUser, ok := user.(models.CurrentUser)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	newToken, err := midware.GenerateNewToken(currentUser)
	if err != nil {
		//
		return
	}
	// 生成新token
	utils.ResponseSuccess(c, newToken)
	return
	// update := bson.M{"$set": bson.M{}}

}
