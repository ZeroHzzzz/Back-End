package studenthandler

import (
	"hr/app/service"
	"hr/app/utils"
	"hr/configs/models"

	"github.com/gin-gonic/gin"
)

type information struct {
	UserID   string `json:"UserID"`
	Category string `json:"Vategory"`
	Content  string `json:"Content"`
}

func FeedbackOAdvice(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var information information
	err := c.ShouldBindJSON(&information)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}
	newFeeback := models.Feedback{
		Category: information.Category,
		UserID:   information.UserID,
		Content:  information.Content,
		Status:   false,
	}
	insertResult := service.InsertOne(c, utils.MongodbName, utils.FeedbackOAdvice, newFeeback)
	utils.ResponseSuccess(c, insertResult.InsertedID) //返回文档的id
}

// func Advice(c *gin.Context) {
// 	c.Header("Content-Type", "application/json")
// 	var information information
// 	err := c.ShouldBindJSON(&information)
// 	if err != nil {
// 		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
// 		c.Abort()
// 		return
// 	}
// 	newFeeback := models.Feedback{
// 		Category: information.Category,
// 		UserID:   information.UserID,
// 		Content:  information.Content,
// 		Status:   false,
// 	}
// 	insertResult := service.InsertOne(c, , "", newFeeback)
// 	utils.ResponseSuccess(c, insertResult.InsertedID) //返回文档的id
// }
