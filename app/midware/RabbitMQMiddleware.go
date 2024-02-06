package midware

import (
	"hr/app/service"

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
		c.Next()
	}
}
