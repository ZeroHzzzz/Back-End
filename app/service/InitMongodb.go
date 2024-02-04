package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MongoDBHost     = "localhost"
	MongoDBPort     = 27017
	MongoDBPassword = "password" // 这里没有配置密码，很有可能出问题
)

// 初始化 MongoDB 客户端
func InitMongoClient() (*mongo.Client, error) {
	// 设置 MongoDB 连接配置
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", MongoDBHost, MongoDBPort)).SetConnectTimeout(10 * time.Second)

	// 连接 MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// 设置最大连接池大小
	clientOptions.SetMaxPoolSize(10)

	// 创建 MongoDB 连接上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return client, nil
}

func CloseMongoClient(client *mongo.Client) {
	if client != nil {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("MongoDB connection closed.")
	}
}
