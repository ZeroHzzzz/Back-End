package counsellorcontroller

import (
	scoredatabase "hr/app/service/square"
	"hr/app/utils"
	"hr/configs/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
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

	list, err = scoredatabase.GetTopicList(getsubmissionlistinformation.Start, getsubmissionlistinformation.End, collection)
	if err != nil {
		//处理逻辑
	}
	utils.ResponseSuccess(c, topiclist)
	return
}
