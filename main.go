package main

import (
	"hr/configs/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// database.Init()
	// 考虑写一个数据库初始化的，比如声明一个交换机
	r := gin.Default()
	router.Init(r)
	err := r.Run(":3000")
	if err != nil {
		log.Fatal("Server start error", err)
	}

}
