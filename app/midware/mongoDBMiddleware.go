package midware

import (
	"context"
	"fmt"
	"hr/app/utils"
	configs "hr/configs/config"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoClientMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取配置
		mongoDBurl := configs.Config.GetString("MongoDB.url")
		mongoDBuser := configs.Config.GetString("MongoDB.user")
		mongoDBpassword := configs.Config.GetString("MongoDB.password")

		// 初始化 MongoDB 客户端
		// 设置 MongoDB 连接配置
		clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", mongoDBuser, mongoDBpassword, mongoDBurl)).SetConnectTimeout(10 * time.Second)
		// 连接 MongoDB
		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			c.Error(utils.GetError(utils.MONGODB_INIT_ERROR, err.Error()))
			c.Abort()
			return
		}
		defer CloseMongoClient(client)
		// 设置最大连接池大小
		clientOptions.SetMaxPoolSize(10)

		// 创建 MongoDB 连接上下文
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel() // 这里会出问题的，应该放在handler里面而不是放在函数内部

		// 检查连接
		err = client.Ping(ctx, nil)
		if err != nil {
			c.Error(utils.GetError(utils.CONNECT_ERROR, err.Error()))
			c.Abort()
			return
		}

		fmt.Println("Connected to MongoDB!")

		// 将 MongoDB 客户端添加到请求的上下文中
		ctx = context.WithValue(c.Request.Context(), "mongoClient", client)

		// 设置上下文为新的带有 MongoDB 客户端的上下文
		c.Request = c.Request.WithContext(ctx)

		// 调用下一个处理程序，最后回到中间件
		c.Next()
	}
}

func CloseMongoClient(client *mongo.Client) {
	if client != nil {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Fatal(err) // 记录错误，这里就不中断了
		}
		fmt.Println("MongoDB connection closed.")
	}
}
