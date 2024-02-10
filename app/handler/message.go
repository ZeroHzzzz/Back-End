package handler

import (
	"hr/app/service"

	"github.com/gin-gonic/gin"
)

func WebSocketConnection(c *gin.Context) {
	userId := c.Param("userID")
	service.HandleWebSocketConnection(c, userId)
}
