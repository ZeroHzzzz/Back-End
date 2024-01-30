package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// func ConnectToMongoDB() (*mongo.Client, error) {
// 	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// 设置连接池参数
// 	clientOptions.SetPoolMonitor(&event.PoolMonitor{})

// 	// 创建一个新的 MongoDB 客户端
// 	client, err := mongo.NewClient(clientOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}

// 	// 连接 MongoDB
// 	err = client.Connect(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}

// 	// 检查连接是否可用
// 	err = client.Ping(ctx, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, err
// 	}

//		fmt.Println("Connected to MongoDB!")
//		return client, nil
//	}

const (
	MongoDBHost     = "localhost"
	MongoDBPort     = 27017
	MongoDBDatabase = "your_database_name"
)

func initMongoClient() error {
	// 设置 MongoDB 连接配置
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", MongoDBHost, MongoDBPort))

	// 连接 MongoDB
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return err
	}

	// 设置最大连接池大小
	clientOptions.SetMaxPoolSize(10)

	// 创建 MongoDB 连接上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 连接 MongoDB
	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}
}

// 关闭 MongoDB 连接
func closeMongoClient(client *mongo.Client) {
	if client != nil {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("MongoDB connection closed.")
	}
}
