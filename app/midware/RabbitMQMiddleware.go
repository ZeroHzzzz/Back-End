package midware

import (
	"hr/app/service"
	"hr/app/utils"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type RabbitMQMiddleware struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func GetRabbitMQMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rabbitmqMiddleware := service.Initrmq(c)
		if rabbitmqMiddleware == nil {
			// 处理初始化失败的情况
			c.Error(utils.RMQ_INIT_ERROR)
			c.Abort()
			return
		}
		defer service.Closermq(rabbitmqMiddleware)
		// 将 rabbitmqMiddleware 设置到上下文中
		c.Set("RabbitMQMiddleware", rabbitmqMiddleware)

		// currentUser
		currentUser := service.GetCurrentUser(c)
		service.DeclareQueue(c, currentUser.UserID) // 声明队列
		// 声明交换机
		// 用户信息交换机
		service.DeclareExchange(c, utils.UserExchange, "direct")
		// 全局信息交换机
		service.DeclareExchange(c, utils.GlobalExchange, "fanout")

		// 绑定系统交换机
		service.BindQueue(c, currentUser.UserID, "", utils.GlobalExchange)
		// 绑定用户交换机
		service.BindQueue(c, currentUser.UserID, currentUser.UserID, utils.UserExchange)
		c.Next()
	}
}
