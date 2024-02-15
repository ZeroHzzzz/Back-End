package counsellorhandler

import (
	"context"
	"fmt"
	"hr/app/service"
	"hr/app/utils"
	"hr/configs/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type auditOneInformation struct {
	AuthorID     string `json:"AuthorID"` // 这里指的是提交人的id
	AcademicYear string `json:"AcademicYear"`
	ItemName     string `json:"ItemName"`
	ItemValue    int64  `json:"ItemValue"`
	Msg          string `json:"Msg"`
	Status       bool   `json:"Status"`
	Cause        string `json:"Cause"`
	Advice       string `json:"Advice"`
}

func AuditOne(c *gin.Context) {
	// 审批单个申报
	c.Header("Content-Type", "application/json")

	// 获取用户信息
	currentUser := service.GetCurrentUser(c)

	var information auditOneInformation
	err := c.ShouldBindJSON(&information)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}

	submissionID := c.Query("SubmissionID")
	objectID, err := primitive.ObjectIDFromHex(submissionID)
	if err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return
	}
	// 更新状态
	filter := bson.M{
		"_id": objectID,
	}
	modified := bson.M{
		"$set": bson.M{
			"AduiterID": currentUser.UserID,
			"Status":    information.Status,
			"Cause":     information.Cause,
			"Advice":    information.Advice,
		},
	}
	service.UpdateOne(c, utils.MongodbName, utils.Submission, filter, modified)
	// 加分
	if information.Status {
		// filter = bson.M{
		// 	"_id":          information.AuthorID,
		// 	"academicYear": information.AcademicYear,
		// 	"itemName":     information.ItemName,
		// 	"mark":         information.ItemValue,
		// 	"msg":          information.Msg,
		// }
		newScore := models.Score{
			UserID:       information.AuthorID,
			AcademicYear: information.AcademicYear,
			ItemName:     information.ItemName,
			Mark:         information.ItemValue,
			Msg:          information.Msg,
		}
		// 改更新为插入，前面要多传入一个msg说明这个分数的来历，然后还有把加分部分改了，然后这个部分一定要多注意
		_ = service.InsertOne(c, utils.MongodbName, utils.Score, newScore)
	}
	// 新建历史记录
	var msg string
	if information.Status {
		msg = fmt.Sprintf("%s通过了申报表%s", currentUser.UserID, submissionID)
	} else {
		msg = fmt.Sprintf("%s驳回了申报表%s", currentUser.UserID, submissionID)
	}
	// TODO：这里后续可以考虑加一个通过超链接直接点击获取申报表信息的，但是好像有接口可以复用
	newHistory := models.SubmitHistory{
		SubmissionID: submissionID,
		AuditorID:    currentUser.UserID,
		Message:      msg,
		CreateAt:     time.Now().Unix(),
	}
	service.InsertOne(c, utils.MongodbName, utils.SubmitHistory, newHistory)

	// 发送通知
	service.PublishMessage(c, utils.UserExchange, information.AuthorID, msg)

	utils.ResponseSuccess(c, nil)
}

type auditManyInformation struct {
	AcademicYear  string   `json:"AcademicYear"`
	SubmissionIDs []string `json:"SubmissionIDs"`
	AuthorIDs     []string `json:"AuthorIDs"`
	Msg           []string `json:"Msg"`
	ItemName      []string `json:"ItemName"`
	ItemValue     []int64  `json:"ItemValue"`
	Status        []bool   `json:"Status"`
	Advice        []string `json:"Advice"`
	Cause         []string `json:"Cause"`
}

func AuditMany(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	// 获取用户信息
	currentUser := service.GetCurrentUser(c)
	var information auditManyInformation
	err := c.ShouldBindJSON(&information)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}
	// 更新申报表状态
	filter := bson.M{
		"_id": bson.M{
			"$in": information.SubmissionIDs,
		},
	}
	modified := bson.M{
		"$set": bson.M{
			"AuditerID": currentUser.UserID,
			"Status":    information.Status,
			"Cause":     information.Cause,
			"Advice":    information.Advice,
		},
	}
	_ = service.UpdateMany(c, utils.MongodbName, utils.Submission, filter, modified)

	// 找出成功的提交，这是为了找出提交人，用来发信和加分
	baseSubmissionHistory := models.SubmitHistory{
		AuditorID: currentUser.UserID,
	}
	var docs []interface{}
	var msg string
	for i, submissionID := range information.SubmissionIDs {
		if information.Status[i] {
			// 加分
			newScore := models.Score{
				UserID:       information.AuthorIDs[i],
				AcademicYear: information.AcademicYear,
				ItemName:     information.ItemName[i],
				Mark:         information.ItemValue[i],
				Msg:          information.Msg[i],
			}
			_ = service.InsertOne(c, utils.MongodbName, utils.Score, newScore)
			msg = fmt.Sprintf("%s通过了申报表%s", currentUser.UserID, submissionID)
		} else {
			msg = fmt.Sprintf("%s驳回了申报表%s", currentUser.UserID, submissionID)
		}
		// 加入历史
		doc := baseSubmissionHistory
		doc.SubmissionID = submissionID
		doc.Message = msg
		doc.CreateAt = time.Now().Unix()
		docs = append(docs, doc)

		// 发信
		service.PublishMessage(c, utils.UserExchange, information.AuthorIDs[i], msg)
	}
	_ = service.InsertMany(c, utils.MongodbName, utils.SubmitHistory, docs)
	utils.ResponseSuccess(c, nil)
}

func GetAuditHistory(c *gin.Context) {
	c.Header("Content-Type", "application/json")
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

	filter := bson.D{}
	options := options.Find().SetSort(bson.D{{Key: "CreateAt", Value: -1}}).SetSkip((int64(page - 1)) * int64(limit)).SetLimit(int64(limit))
	result := service.Find(c, utils.MongodbName, utils.SubmitHistory, filter, options)

	var list []models.SubmitHistory
	if err := result.All(context.TODO(), &list); err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return
	}
	utils.ResponseSuccess(c, list)
}
