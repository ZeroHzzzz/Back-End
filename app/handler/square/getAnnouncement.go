package squarehandler

import (
	"context"
	counsellorhandler "hr/app/handler/counsellor"
	"hr/app/service"
	"hr/app/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAnnouncement(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	filter := bson.M{}
	options := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(5)
	// 找五条最新的
	var list []counsellorhandler.Announcement
	cursor := service.Find(c, utils.MongodbName, utils.Announcement, filter, options)
	if err := cursor.All(context.TODO(), &list); err != nil {
		c.Error(utils.GetError(utils.DECODE_ERROR, err.Error()))
		c.Abort()
		return
	}
	utils.ResponseSuccess(c, list)
	return
}
