package midware

import (
	"context"
	"hr/app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func mongoClientMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 初始化 MongoDB 客户端
		client, err := utils.InitMongoClient()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize MongoDB client"})
			c.Abort()
			return
		}

		// 将 MongoDB 客户端添加到请求的上下文中
		ctx := context.WithValue(c.Request.Context(), "mongoClient", client)

		// 设置上下文为新的带有 MongoDB 客户端的上下文
		c.Request = c.Request.WithContext(ctx)

		// 调用下一个处理程序，最后回到中间件
		c.Next()

		utils.CloseMongoClient(client)
	}
}
