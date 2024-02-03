package studentcontroller

import (
	"hr/app/midware"
	"hr/app/service"
	"hr/app/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type modifiedprofileInformation struct {
	PassWord    string `json:"passWord"`
	NewPassword string `json:"newPassword"`
}

func ModifiedProfileHandler(c *gin.Context) {
	//修改密码后token行为可能需要深入考虑
	c.Header("Content-Type", "application/json")
	var modifiedprofileinformation modifiedprofileInformation
	err := c.ShouldBindJSON(&modifiedprofileinformation)
	if err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}

	userid := c.Param("userId")

	filter := bson.M{
		"userId":   userid,
		"passWord": modifiedprofileinformation.PassWord,
	}
	modified := bson.M{
		"$set": bson.M{
			"passWord": modifiedprofileinformation.NewPassword,
		},
	}
	// 修改之后的文档
	_ = service.UpdateOne(c, "", "", filter, modified)
	if err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}
	currentUser := service.GetCurrentUser(c)
	newToken, err := midware.GenerateToken(currentUser)
	if err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}
	// 生成新token
	utils.ResponseSuccess(c, newToken)
}
