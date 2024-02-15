package studenthandler

import (
	"hr/app/service"
	"hr/app/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type modifiedprofileInformation struct {
	PassWord    string `json:"PassWord"`
	NewPassword string `json:"NewPassword"`
}

func ModifiedProfileHandler(c *gin.Context) {
	//修改密码后token行为可能需要深入考虑
	c.Header("Content-Type", "application/json")
	var modifiedprofileinformation modifiedprofileInformation
	err := c.ShouldBindJSON(&modifiedprofileinformation)
	if err != nil {
		c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
		c.Abort()
		return
	}

	userid := c.Param("UserID")

	filter := bson.M{
		"_id":      userid,
		"PassWord": modifiedprofileinformation.PassWord,
	}
	modified := bson.M{
		"$set": bson.M{
			"PassWord": modifiedprofileinformation.NewPassword,
		},
	}
	// 修改之后的文档
	_ = service.UpdateOne(c, utils.MongodbName, utils.Student, filter, modified)
	currentUser := service.GetCurrentUser(c)

	// newToken, err := midware.GenerateToken(currentUser)
	// if err != nil {
	// 	c.Error(utils.GetError(utils.PARAM_ERROR, err.Error()))
	// 	c.Abort()
	// }
	// 发信
	service.PublishMessage(c, utils.UserExchange, currentUser.UserID, utils.ModifiedProfile)
	utils.ResponseSuccess(c, nil)
}
