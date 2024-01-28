package router

import (
	"hr/app/midware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	// 跨域
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	const prelogin = "/login"
	login := r.Group(prelogin)
	{
		login.GET("/student")
		login.GET("/counsellor")
	}
	const preStudent = "/student"
	student := r.Group(preStudent)
	{
		student.Use(func(c *gin.Context) {
			midware.AuthenticateMiddleware(c, "admin", "student")
			if !c.IsAborted() {
				c.Next()
			}
		})
		student.PUT("/profile")
	}
}
