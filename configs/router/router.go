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

	// test
	test := r.Group("/test", midware.ErrorHandler(), midware.MongoClientMiddleware())
	{
		test.GET("", handler.Test)
	}
	login := r.Group("/login", midware.ErrorHandler(), midware.MongoClientMiddleware())
	{
		login.GET("/student", handler.LoginHandler_Student)
		login.GET("/counsellor", handler.LoginHandler_Counsellor)
	}
	api := r.Group("/api", midware.ErrorHandler(), midware.MongoClientMiddleware())
	{
		api.GET("/ws/:userId", handler.WebSocketConnection)
		student := api.Group("/student", midware.JWTAuthMiddleware("Counsellor", "Student"), midware.GetRabbitMQMiddleware(), midware.RedisClientMiddleware())
		{
			student.PUT("/profile/:userId", studenthandler.ModifiedProfileHandler)
			student.POST("/feedbackOradvice", studenthandler.FeedbackOAdvice)
			student.GET("/score", studenthandler.GetConcreteSorce)
			submit := student.Group("/submit")
			{
				submit.POST("/:userId", studenthandler.Submission)
				submit.GET("/status/:formID", studenthandler.GetSubmissionStatus)
				submit.GET("/list", studenthandler.GetSubmissionStatus)
			}
		}
		counsellor := api.Group("/counsellor", midware.JWTAuthMiddleware("Counsellor"), midware.GetRabbitMQMiddleware(), midware.RedisClientMiddleware())
		{
			counsellor.POST("/:counsellorId/cause", counsellorhandler.AddCause)
			counsellor.GET("/:counsellorId/cause", counsellorhandler.GetCause)
			counsellor.POST("/access-time", counsellorhandler.SetAccessTimeHandler)
			counsellor.POST("/setannouncement", counsellorhandler.SetAnnouncement)
			audit := counsellor.Group("/audit")
			{
				audit.GET("/list", counsellorhandler.GetSubmissionList)
				audit.PUT("/review/single", counsellorhandler.AuditOne)
				audit.PUT("/review/bulk", counsellorhandler.AuditMany)
				audit.GET("/history", counsellorhandler.GetAuditHistory)
				// audit.PUT("/remake/:submissionId", counsellorhandler.)
			}
			information := counsellor.Group("/information")
			{
				information.PUT("/correct/:userID", counsellorhandler.CorrectGrade)
				information.POST("/bulk-import/student", counsellorhandler.ImportStudent)
				information.POST("/bulk-import/counsellor", counsellorhandler.ImportCounsellor)
				information.POST("/bulk-import/mark", counsellorhandler.ImportMark)
				information.GET("/student", counsellorhandler.GetStudentInformation)
			}
		}
		square := api.Group("/square", midware.JWTAuthMiddleware("Counsellor", "Student"), midware.GetRabbitMQMiddleware(), midware.RedisClientMiddleware())
		{
			square.GET("/annoucement", squarehandler.GetAnnouncement)
			topic := square.Group("/topic")
			{
				topic.POST("/new", squarehandler.NewTopic)
				topic.PUT("", squarehandler.ModifiedTopic)
				topic.GET("/list", squarehandler.GetTopicList)
				topic.GET("", squarehandler.GetTopic)
				topic.POST("/replies", squarehandler.NewReply)
				topic.GET("/replies", squarehandler.GetReply)
				topic.GET("/views&likes", squarehandler.GetViewsAndlikes)
				topic.PUT("/likes/reply", squarehandler.LikeReply)
				topic.PUT("/likes/topic", squarehandler.LikesTopic)
				topic.DELETE("/delete/topic", squarehandler.DeleteTopic)
				topic.DELETE("/delete/reply", squarehandler.DeleteReply)
			}
		}
	}
}
