package handler

import (
	"hr/app/service"

	"github.com/gin-gonic/gin"
)

func handleWebSocketConnection(c *gin.Context) {
	userId := c.Param("userID")
	service.DeclareExchange(c, "") // 填入交换机的名字
	service.HandleWebSocketConnection(c, userId)
}
