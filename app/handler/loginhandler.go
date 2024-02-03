package handler

import (
	"hr/app/service"
	"hr/app/utils"
	"hr/configs/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func LoginHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	userid := c.PostForm("userId")
	password := c.PostForm("passWord")
	var user models.Student
	// err := database.DB.Where("userId=?", "passWord=?", userid, password).First(&user).Error

	filter := bson.M{
		"userId":   userid,
		"passWord": password,
	}

	err := service.FindOne(c, "", "", filter).Decode(&user)
	if err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		c.Abort()
	}

	currentUser := models.CurrentUser{
		UserId:     userid,
		UserName:   user.UserName,
		Grade:      user.Grade,
		Profession: user.Profession,
	}
	c.Set("CurrentUser", currentUser) //将用户信息储存到上下文
	utils.ResponseSuccess(c, currentUser)
}
