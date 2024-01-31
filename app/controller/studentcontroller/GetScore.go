package studentcontroller

import (
	"hr/app/service/student"
	"hr/app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetSorceHandler(c *gin.Context) {
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

	rezult, err := student.GetYearScore(userId, collection)
	if err != nil {
		// ...
		return
	}
	utils.ResponseSuccess(c, rezult)
	return

}
