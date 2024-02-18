package midware

import (
	"hr/app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		lastError := c.Errors.Last()
		if lastError != nil {
			if c.IsAborted() {
				err := lastError.Err
				// 若是自定义的错误则将code、msg返回
				if myErr, ok := err.(*utils.MyError); ok {
					c.JSON(http.StatusOK, gin.H{
						"code": myErr.Code,
						"msg":  myErr.Msg,
						"data": myErr.Data,
					})
					return
				}
				// 若非自定义错误则返回详细错误信息err.Error()
				// 比如save session出错时设置的err
				c.JSON(http.StatusOK, gin.H{
					"code": 500,
					"msg":  "服务器异常",
					"data": err.Error(),
				})
				return
			}
		}
	}
}
