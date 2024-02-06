package counsellor

import (
	"hr/app/service"
	"hr/app/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type announcement struct {
	Content string `json:"content"`
}

type Announcement struct {
	AutherID string    `bson:"autherID"`
	Content  string    `json:"content"`
	CreateAt time.Time `bson:"createdAt"`
}

func SetAnnouncement(c *gin.Context) {
	var information announcement
	err := c.ShouldBindJSON(&information)
	if err != nil {
		c.Error(utils.GetError(utils.VALID_ERROR, err.Error()))
		return
	}
	currentUser := service.GetCurrentUser(c)
	newAnnouncement := Announcement{
		AutherID: currentUser.UserId,
		Content:  information.Content,
	}
	_ = service.InsertOne(c, "", "", newAnnouncement)
	service.PublishMessage(c, utils.GlobalExchange, "", currentUser.UserName+utils.Announcement+": "+information.Content) // 发布信息 用扇out交换机
	utils.ResponseSuccess(c, nil)
	return
}
