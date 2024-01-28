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
			midware.AuthenticateMiddleware(c, "counsellor", "student")
			if !c.IsAborted() {
				c.Next()
			}
		})

		student.PUT("/profile/:userId")
		student.GET("/profile/:userId")

		student.POST("/:userId/form/submit")
		student.GET("/:userId/form/search")

		student.POST("/feedback")
		student.POST("/recommend")

		student.GET("/score/:userId/:academicYear")
	}

	const preCounsellor = "/counsellor"
	counsellor := r.Group(preCounsellor)
	{
		counsellor.Use(func(c *gin.Context) {
			midware.AuthenticateMiddleware(c, "counsellor")
			if !c.IsAborted() {
				c.Next()
			}
		})

		counsellor.GET("/audit/pending")
		counsellor.GET("/audit/:formId")
		counsellor.PUT("/audit/:formId/review")
		counsellor.PUT("/audit/review")
		counsellor.GET("/audit/:formId/history")
		counsellor.DELETE("/audit/:formId/revoke")

		counsellor.PUT("/grade/personal-import/:userId")
		counsellor.PUT("/grade/bulk-import")

		counsellor.POST("/:counsellorId/cause")
		counsellor.POST("/setddl")
	}

	const preSquare = "square"
	square := r.Group(preSquare)
	{
		square.Use(func(c *gin.Context) {
			midware.AuthenticateMiddleware(c, "counsellor", "student")
			if !c.IsAborted() {
				c.Next()
			}
		})
		square.POST("/topic/create")
		square.GET("/topic")
		square.GET("/square/topic/:topicId")
		square.POST("/topic/:topicId/replies")
		square.PUT("/topic/:topicId")
		square.DELETE("/topic/:topicId")
	}
}
