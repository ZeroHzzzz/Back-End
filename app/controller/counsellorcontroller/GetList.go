package counsellorcontroller

import (
	"context"
	"hr/app/utils"
	"hr/configs/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type getSubmissionListInformation struct {
	Start int64 `json:"start" binding:"required"`
	End   int64 `json:"end" binding:"required"`
}

func GetTopicList(c *gin.Context) {
	const DatabaseName string = ""
	const CollectionName string = ""

	var getsubmissionlistinformation getSubmissionListInformation
	err := c.ShouldBindJSON(&getsubmissionlistinformation)
	if err != nil {
		utils.ResponseError(c, "Paramter", "ParameterErrorMsg")
		return
	}

	// 获取collection
	var list []models.SubmitInformation
	mongoClient, exists := c.Request.Context().Value("mongoClient").(*mongo.Client)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MongoDB client not found in context"})
		return
	}
	database := mongoClient.Database(DatabaseName)
	collection := database.Collection(CollectionName)
	filter := bson.D{}
	options := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip(getsubmissionlistinformation.Start).SetLimit(getsubmissionlistinformation.End - getsubmissionlistinformation.Start + 1)

	// 执行查询
	cursor, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		// 处理
		return
	}
	if err := cursor.All(context.TODO(), &list); err != nil {
		return
	}

	utils.ResponseSuccess(c, list)
	return
}
