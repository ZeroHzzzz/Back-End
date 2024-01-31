package main

import (
	"hr/configs/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// database.Init()
	r := gin.Default()
	router.Init(r)
	// err := r.Run(":3000")
	// if err != nil {
	// 	log.Fatal("Server start error", err)
	// }

}
