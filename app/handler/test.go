package handler

import (
	"hr/app/service"
	"hr/app/utils"
	"hr/configs/models"

	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	counsellor := models.Counsellor{
		UserId:     "admin",
		UserName:   "Stack",
		PassWord:   "admin",
		Grade:      "大一",
		Profession: "计算机",
	}
	_ = service.InsertOne(c, utils.MongodbName, utils.Counsellor, counsellor)
	utils.ResponseSuccess(c, "dfsdfasfadfa")
}
