package counsellorhandler

import (
	"context"
	"hr/app/service"
	"hr/app/utils"
	"hr/configs/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type getSubmissionListInformation struct {
	Profession string `json:"Profession"`
	Grade      string `json:"Grade"` //年级
	Class      string `json:"Class"`
}

func GetSubmissionList(c *gin.Context) {
	var getsubmissionlistinformation getSubmissionListInformation
	pageParam := c.Query("Page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}
	limitParam := c.Query("Limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}
	err = c.ShouldBindJSON(&getsubmissionlistinformation)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}

	// 获取collection
	var list []models.SubmitInformation

	// 获取未审核表单
	// filter := bson.M{
	// 	"class":      getsubmissionlistinformation.Class,
	// 	"profession": getsubmissionlistinformation.Profession,
	// 	"grade":      getsubmissionlistinformation.Grade,
	// 	"status":     false,
	// }
	filter := bson.M{"Status": false}
	if getsubmissionlistinformation.Class != "" {
		filter["Class"] = getsubmissionlistinformation.Class
	}
	if getsubmissionlistinformation.Profession != "" {
		filter["Profession"] = getsubmissionlistinformation.Profession
	}
	if getsubmissionlistinformation.Grade != "" {
		filter["Grade"] = getsubmissionlistinformation.Grade
	}
	options := options.Find().SetSort(bson.D{{Key: "CreateAt", Value: -1}}).SetSkip(int64((page - 1) * limit)).SetLimit(int64(limit))

	// 执行查询
	cursor := service.Find(c, utils.MongodbName, utils.Submission, filter, options)
	if err := cursor.All(context.TODO(), &list); err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return
	}
	utils.ResponseSuccess(c, list)
}

func GetSubmission(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	submissionID := c.Param("SubmissionID")

	filter := bson.M{"SubmissionID": submissionID}
	result := service.FindOne(c, utils.MongodbName, utils.Submission, filter)

	var submission models.SubmitInformation
	err := result.Decode(&submission)
	if err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return
	}
	utils.ResponseSuccess(c, submission)
}
