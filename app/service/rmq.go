package service

import (
	"fmt"
	"hr/app/utils"
	"hr/configs/models"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

const (
	rmqHost     = "localhost"
	rmqPort     = 6379
	rmqPassword = ""
)

func Initrmq(c *gin.Context) *models.RabbitMQMiddleware {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@localhost:%d/", rmqPort))
	if err != nil {
		c.Error(err) // TODO:
		c.Abort()    // TODO
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		c.Error(err)
		c.Abort()
	}
	// 声明交换机
	// 用户信息交换机
	DeclareExchange(c, utils.UserExchange, "direct")
	// 全局信息交换机
	DeclareExchange(c, utils.GlobalExchange, "fanout")

	return &models.RabbitMQMiddleware{
		Connection: conn,
		Channel:    ch,
	}

}

func Closermq(r *models.RabbitMQMiddleware) {
	if r != nil {
		r.Channel.Close()
		r.Connection.Close()
	}
}
