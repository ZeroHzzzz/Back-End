package counsellorhandler

import (
	"context"
	"hr/app/service"
	"hr/app/utils"
	"hr/configs/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type addCause struct {
	Msg string `json:"msg"`
}

func AddCause(c *gin.Context) {
	userId := c.Param("counsellorId")
	var information addCause
	if err := c.BindJSON(&information); err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}
	newCause := models.Cause{
		UserId: userId,
		Msg:    information.Msg,
	}
	_ = service.InsertOne(c, utils.MongodbName, utils.Cause, newCause)
	utils.ResponseSuccess(c, nil)
}

func GetCause(c *gin.Context) {
	userId := c.Param("counsellorId")
	filter := bson.M{
		"_id": userId,
	}
	var list []models.Cause
	cursor := service.Find(c, utils.MongodbName, utils.Cause, filter)
	if err := cursor.All(context.TODO(), &list); err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return
	}
	utils.ResponseSuccess(c, list)
}
