package router

import (
	"hr/app/handler"
	counsellorhandler "hr/app/handler/counsellor"
	squarehandler "hr/app/handler/square"
	"hr/app/handler/studenthandler"
	"hr/app/midware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 记得应该进行数据库连接，然后再进行鉴权
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
	login := r.Group("login", midware.ErrorHandler(), midware.MongoClientMiddleware())
	{
		login.GET("/student", handler.LoginHandler_Student)
		login.GET("/counseller", handler.LoginHandler_Counsellor)
	}
	api := r.Group("/api", midware.ErrorHandler(), midware.MongoClientMiddleware(), midware.JWTAuthMiddleware("Student", "Counsellor"))
	{
		api.GET("/ws/:userID", handler.WebSocketConnection)
		api.GET("/submit/history", counsellorhandler.GetAuditHistory)
		student := api.Group("/student", midware.GetRabbitMQMiddleware(), midware.RedisClientMiddleware())
		{
			student.PUT("/profile/:userId", studenthandler.ModifiedProfileHandler)
			student.POST("/feedbackOradvice", studenthandler.FeedbackOAdvice)
			student.GET("/score", studenthandler.GetConcreteSorce)
			submit := student.Group("/submit")
			{
				submit.POST("/submission", studenthandler.Submission)
				submit.GET("/status/:formID", studenthandler.GetSubmissionStatus)
				submit.GET("/list")
			}
		}
		counsellor := api.Group("/counsellor", midware.GetRabbitMQMiddleware(), midware.RedisClientMiddleware())
		{
			counsellor.POST("/:counsellorId/cause", counsellorhandler.AddCause)
			counsellor.GET("/:counsellorId/cause", counsellorhandler.GetCause)
			counsellor.POST("/access-time", counsellorhandler.SetAccessTimeHandler)
			counsellor.POST("/setannouncement", counsellorhandler.SetAnnouncement)
			audit := counsellor.Group("/audit")
			{
				audit.GET("/list", counsellorhandler.GetSubmissionList)
				audit.PUT("/review/:submission", counsellorhandler.AuditOne)
				audit.PUT("review", counsellorhandler.AuditMany)
				// audit.PUT("/remake/:submissionId", counsellorhandler.)
			}
			grade := counsellor.Group("/grade")
			{
				grade.PUT("/correct/:userID", counsellorhandler.CorrectGrade)
				grade.POST("/bulk-import", counsellorhandler.ImportStudentInformation)
			}
		}
		square := api.Group("/square")
		{
			square.GET("/annnoucement", squarehandler.GetAnnouncement)
			topic := square.Group("/topic")
			{
				topic.POST("/new", squarehandler.NewTopic)
				topic.PUT("/:topicID", squarehandler.ModifiedTopic)
				topic.GET("/list", squarehandler.GetTopicList)
				topic.GET("/:topicID", squarehandler.GetTopic)
				topic.POST("/:topicID/replies", squarehandler.NewReply)
				topic.GET("/:topicID/views&likes0", squarehandler.GetViewsAndlikes)
				topic.PUT("/reply/:replyID/likes", squarehandler.LikeReply)
				topic.PUT("/:topicID/likes", squarehandler.LikesTopic)
				topic.DELETE("/:topicID", squarehandler.DeleteTopic)
				topic.DELETE("/reply/:replyID", squarehandler.DeleteReply)
			}
		}
	}
}
