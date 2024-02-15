package service

import (
	"fmt"
	"hr/app/utils"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type UserConnection struct {
	Conn   *websocket.Conn
	UserID string
}

var connections = make(map[string]*UserConnection)
var mu sync.Mutex

func HandleWebSocketConnection(c *gin.Context, userID string) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	userConn := &UserConnection{Conn: conn, UserID: userID}

	// 将用户连接添加到连接池
	mu.Lock()
	connections[userID] = userConn
	mu.Unlock()

	log.Printf("User %s connected\n", userID)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			// 处理连接关闭或出错的情况
			log.Printf("User %s disconnected\n", userID)
			utils.ResponseError(c, err.Error())
			// 将用户连接从连接池中移除
			mu.Lock()
			delete(connections, userID)
			mu.Unlock()

			break
		}
		// 声明一个消费者并进行监听
		msgs := ConsumeMessage(c, userID)
		r := GetRabbitMQMiddle(c)

		go func() {
			for msg := range msgs {
				SendMessageToClient(c, userID, msg.Body)
				err := r.Channel.Ack(msg.DeliveryTag, false)
				if err != nil {
					msg := fmt.Sprintf("Failed to acknowledge message: %v\n", err)
					log.Printf(msg)
					utils.ResponseError(c, msg)
				}
			}
		}()

	}
}

// 发信
func SendMessageToClient(c *gin.Context, userID string, message []byte) {
	mu.Lock()
	defer mu.Unlock()
	userConn, exists := connections[userID]
	if !exists {

		return
	}
	err := userConn.Conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		msg := fmt.Sprintf("Failed to send message to user %s: %v\n", userID, err)
		log.Printf(msg)
		utils.ResponseError(c, msg)
	}
}
