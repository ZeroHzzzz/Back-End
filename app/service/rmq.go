package service

import (
	"fmt"
	"hr/app/utils"
	configs "hr/configs/configuration"
	"hr/configs/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func Initrmq(c *gin.Context) *models.RabbitMQMiddleware {
	// 加载配置
	rabbitMQurl := configs.Config.GetString("rabbitMQ.url")
	rabbitMQuser := configs.Config.GetString("rabbitMQ.user")
	rabbitMQpassword := configs.Config.GetString("rabbitMQ.password")

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s", rabbitMQuser, rabbitMQpassword, rabbitMQurl))
	if err != nil {
		log.Printf("连接 RabbitMQ 失败：%v\n", err)
		return nil
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Printf("创建 RabbitMQ 通道失败：%v\n", err)
		return nil
	}
	return &models.RabbitMQMiddleware{
		Connection: conn,
		Channel:    ch,
	}
}

func Closermq(rmq *models.RabbitMQMiddleware) {
	if rmq == nil {
		return
	}
	if rmq.Channel != nil {
		rmq.Channel.Close()
	}
	if rmq.Connection != nil {
		rmq.Connection.Close()
	}
}

// rmq
func DeclareQueue(c *gin.Context, queueName string) amqp.Queue {
	r := GetRabbitMQMiddle(c)
	if r == nil || r.Channel == nil {
		fmt.Println("没有上下文噢")
		c.Error(utils.GetError(utils.RMQ_INIT_ERROR, "RabbitMQ client or channel is nil"))
		c.Abort()
		return amqp.Queue{} // 返回一个零值队列，或者根据实际情况返回适当的值
	}

	q, err := r.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		true,      // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Printf("声明队列失败：%v\n", err)
		c.Error(utils.GetError(utils.RMQ_DECLARE_QUEUE_ERROR, err.Error()))
		c.Abort()
		return amqp.Queue{} // 返回一个零值队列，或者根据实际情况返回适当的值
	}
	return q
}

func DeclareExchange(c *gin.Context, exchangeName, kind string) {
	// 声明交换机
	r := GetRabbitMQMiddle(c)
	if r == nil || r.Channel == nil {
		c.Error(utils.GetError(utils.RMQ_INIT_ERROR, "RabbitMQ client or channel is nil"))
		c.Abort()
		return
	}

	err := r.Channel.ExchangeDeclare(
		exchangeName,
		kind,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("声明交换机失败：%v\n", err)
		c.Error(utils.GetError(utils.RMQ_DECLARE_EXCHANGE_ERROR, err.Error()))
		c.Abort()
		return
	}
}

func BindQueue(c *gin.Context, queueName, routingKey, exchangeName string) {
	r := GetRabbitMQMiddle(c)
	if r == nil || r.Channel == nil {
		c.Error(utils.GetError(utils.RMQ_INIT_ERROR, "RabbitMQ client or channel is nil"))
		c.Abort()
		return
	}

	err := r.Channel.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange name
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Printf("交换机绑定队列失败：%v\n", err)
		c.Error(utils.GetError(utils.RMQ_BIND_QUEUE_ERROR, err.Error()))
		c.Abort()
	}
}

func ConsumeMessage(c *gin.Context, queueName string) <-chan amqp.Delivery {
	// 声明一个消费者，并返回一个接收（receive）操作符用于从通道中接收数据的表达式
	r := GetRabbitMQMiddle(c)
	msgs, err := r.Channel.Consume(
		queueName, // queueName
		"",        // consumer
		false,     //autoAck 自动确认已读，false为手动
		false,     // exclusive 独占
		true,      // 接收自己的信息
		true,      // 不等待服务器响应，false表示等待
		nil,       // 其他参数
	)
	if err != nil {
		log.Printf("消费信息失败：%v\n", err)
		c.Error(utils.GetError(utils.RMQ_DECLARE_CONSUMER_ERROR, err.Error()))
		c.Abort()
		return nil
	}
	return msgs
}

// for msg := range messages {
// 	// 将消息发送到前端
// 	sendMessageToClient(msg.Body)
// }

func PublishMessage(c *gin.Context, exchangeName, queueName, message string) {
	r := GetRabbitMQMiddle(c)
	err := r.Channel.Publish(
		exchangeName, // exchange
		queueName,    // routing key (队列名即为路由键) ,如果为空就是发布到全部队列
		false,        // mandatory
		false,        //
		//  immediate 参数为 false 时（默认值），如果消息无法被立即路由到队列，消息将会被存储在队列中等待消费者接收
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(message),
			DeliveryMode: amqp.Persistent, // 持久化
		},
	)
	if err != nil {
		log.Printf("发送信息失败：%v\n", err)
		c.Error(utils.GetError(utils.RMQ_PUBLISH_ERROR, err.Error()))
		c.Abort()
		return
	}
}
