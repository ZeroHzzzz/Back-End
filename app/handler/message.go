package handler

import (
	"hr/app/service"

	"github.com/gin-gonic/gin"
)

func WebSocketConnection(c *gin.Context) {
	userID := c.Param("UserID")
	service.HandleWebSocketConnection(c, userID)
}
