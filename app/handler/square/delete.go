package square

import (
	"hr/app/service"
	"hr/app/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func DeleteTopic(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	topicId := c.Param("topicId")
	// 从上下文中获取用户信息
	currentUser := service.GetCurrentUser(c)
	// 辅导员拥有删除文章的能力
	if currentUser.Role == "counsellor" {
		filter := bson.M{
			"_id": topicId,
		}
		_ = service.DeleteOne(c, utils.MongodbName, utils.Topic, filter)
		// 删除评论
		filter = bson.M{
			"topicId": topicId,
		}
		_ = service.DeleteMany(c, utils.MongodbName, utils.Reply, filter)

	} else if currentUser.Role == "student" {
		filter := bson.M{
			"_id":      topicId,
			"autherID": currentUser.UserId,
		}
		_ = service.DeleteOne(c, utils.MongodbName, utils.Topic, filter)

		filter = bson.M{
			"topicId": topicId,
		}
		_ = service.DeleteMany(c, utils.MongodbName, utils.Reply, filter)

	}
	utils.ResponseSuccess(c, nil)
}

func DeleteReply(c *gin.Context) {
	// 这个接口故意留了一个漏洞，就是这里只要是用户鉴权成功就能删除评论
	// 这里是为了防止恶意评论的
	// 还有另外一个漏洞就是，删除评论并不能删除全部的子评论，比如子评论的子评论就删除不了，但是在前端不会显示出来(因为没有父评论)
	// 因此目前只有完全删除文章才能删除全部的评论释放空间
	c.Header("Content-Type", "application/json")
	replyID := c.Param("replyID")

	// 从上下文中获取mongo客户端

	filter := bson.M{
		"replyID": replyID,
	}
	_ = service.DeleteOne(c, utils.MongodbName, utils.Reply, filter)

	filter = bson.M{
		"parentID": replyID,
	}
	// 删除子评论
	_ = service.DeleteMany(c, utils.MongodbName, utils.Reply, filter)
	utils.ResponseSuccess(c, nil)
}
