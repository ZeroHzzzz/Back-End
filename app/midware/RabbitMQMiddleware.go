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
		defer service.Closermq(rabbitmqMiddleware)
		c.Set("RabbitMQMiddleware", rabbitmqMiddleware)
		// currentUser
		currentUser := service.GetCurrentUser(c)
		_ = service.DeclareQueue(c, currentUser.UserId) // 声明队列
		// 这里可以不用把他写到上下文里面，因为用户id就是队列名
		// 绑定系统交换机
		service.BindQueue(c, currentUser.UserId, "", utils.GlobalExchange)
		// 绑定用户交换机
		service.BindQueue(c, currentUser.UserId, currentUser.UserId, utils.UserExchange)
		c.Next()
	}
}
