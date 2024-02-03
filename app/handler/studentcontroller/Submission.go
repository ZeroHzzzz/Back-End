package studentcontroller

import (
	"context"
	"hr/app/service"
	"hr/app/utils"
	"hr/configs/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

const savePath = ""

func SubmitHandler(c *gin.Context) {
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
	currentUser := service.GetCurrentUser(c)
	if currentUser.UserId != userId {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	newSubmission := models.SubmitInformation{
		CurrentUser:  currentUser,
		ItemName:     itemName,
		AcademicYear: academicYear,
		Evidence:     destPaths,
		Status:       false,
	}
	insertResult := service.InsertOne(c, "", "", newSubmission)
	utils.ResponseSuccess(c, insertResult.InsertedID)
}

func GetSubmissionStatus(c *gin.Context) {
	// 从Form数据库中查找，然后返回每个的状态，数据格式应该是一个字典
	c.Header("Content-Type", "application/json")
	userId := c.Param("userId")

	filter := bson.M{
		"userId": userId,
	}

	result := service.Find(c, "", "", filter)
	var forms []models.SubmitInformation
	if err := result.All(context.Background(), &forms); err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}
	utils.ResponseSuccess(c, forms)
}
