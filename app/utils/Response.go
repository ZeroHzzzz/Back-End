package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 需要重写
func Response(c *gin.Context, httpStatusCode int, code int, msg string, data interface{}) {
	c.JSON(httpStatusCode, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}
func ResponseUnauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
}
func ResponseSuccess(c *gin.Context, data interface{}) {
	Response(c, http.StatusOK, 200, "Success", data)
}

func ResponseInternalError(c *gin.Context) {
	Response(c, http.StatusInternalServerError, 500, "InternalServerError", nil)
}

func ResponseError(c *gin.Context, data interface{}) {
	Response(c, http.StatusInternalServerError, 502, "Error", data)
}
