package handler

import (
	"hr/app/midware"
	"hr/app/service"
	"hr/app/utils"
	"hr/configs/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type information struct {
	UserID   string `json:"UserID" binding:"required"`
	PassWord string `json:"PassWord" binding:"required"`
}
type reponse struct {
	CurrentUser models.CurrentUser
	Token       string
}

func LoginHandler_Student(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var information information
	err := c.ShouldBindJSON(&information)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}
	var user models.Student

	filter := bson.M{
		"UserID":   information.UserID,
		"PassWord": information.PassWord,
	}

	result := service.FindOne(c, utils.MongodbName, utils.Student, filter)
	if result == nil {
		c.Error(utils.LOGIN_ERROR)
		c.Abort()
		return
	}
	err = result.Decode(&user)
	if err != nil {
		c.Error(utils.GetError(utils.LOGIN_ERROR, err.Error()))
		c.Abort()
		return
	}

	currentUser := models.CurrentUser{
		UserID:     information.UserID,
		UserName:   user.UserName,
		Grade:      user.Grade,
		Role:       "Student",
		Profession: user.Profession,
	}
	c.Set("CurrentUser", currentUser) //将用户信息储存到上下文
	reponse := reponse{
		CurrentUser: currentUser,
	}
	token, err := midware.GenerateToken(currentUser)
	if err == nil {
		reponse.Token = token
	} else {
		c.Error(utils.GetError(utils.INNER_ERROR, err.Error()))
		c.Abort()
		return
	}
	utils.ResponseSuccess(c, reponse)
}

func LoginHandler_Counsellor(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var information information
	err := c.ShouldBindJSON(&information)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}
	var user models.Counsellor
	filter := bson.M{
		"UserID":   information.UserID,
		"PassWord": information.PassWord,
	}

	result := service.FindOne(c, utils.MongodbName, utils.Counsellor, filter)
	if result == nil {
		c.Error(utils.LOGIN_ERROR)
		c.Abort()
		return
	}
	err = result.Decode(&user)
	if err != nil {
		c.Error(utils.GetError(utils.LOGIN_ERROR, err.Error()))
		c.Abort()
		return
	}

	currentUser := models.CurrentUser{
		UserID:     information.UserID,
		UserName:   user.UserName,
		Grade:      user.Grade,
		Role:       "Counsellor",
		Profession: user.Profession,
	}
	// c.Set("CurrentUser", currentUser) //将用户信息储存到上下文
	reponse := reponse{
		CurrentUser: currentUser,
	}
	token, err := midware.GenerateToken(currentUser)
	if err == nil {
		reponse.Token = token
	} else {
		c.Error(utils.GetError(utils.INNER_ERROR, err.Error()))
		c.Abort()
		return
	}
	utils.ResponseSuccess(c, reponse)
}
