package controller

import (
	"hr/app/utils"
	database "hr/configs/database/mySQL"
	"hr/configs/models"

	"github.com/gin-gonic/gin"
)

type CurrentUser struct {
	UserId     string
	userName   string
	grade      int
	profession string
}

func loginHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	userid := c.PostForm("userId")
	password := c.PostForm("passWord")
	var user models.Student
	err := database.DB.Where("userId=?", "passWord=?", userid, password).First(&user).Error
	if err != nil {
		//处理逻辑
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
