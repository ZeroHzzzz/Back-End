package studentcontroller

import (
	"hr/app/service"
	"hr/app/utils"
	"hr/configs/models"

	"github.com/gin-gonic/gin"
)

type information struct {
	UserId   string `json:"userId"`
	Category string `json:"category"`
	Content  string `json:"content"`
}

func Feedback(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var information information
	err := c.ShouldBindJSON(&information)
	if err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}
	newFeeback := models.Feedback{
		Category: information.Category,
		UserId:   information.UserId,
		Content:  information.Content,
		Status:   false,
	}
	insertResult := service.InsertOne(c, "", "", newFeeback)
	utils.ResponseSuccess(c, insertResult.InsertedID) //返回文档的id
}

func Advice(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var information information
	err := c.ShouldBindJSON(&information)
	if err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}
	newFeeback := models.Feedback{
		Category: information.Category,
		UserId:   information.UserId,
		Content:  information.Content,
		Status:   false,
	}
	insertResult := service.InsertOne(c, "", "", newFeeback)
	utils.ResponseSuccess(c, insertResult.InsertedID) //返回文档的id
}
