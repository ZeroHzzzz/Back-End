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

// Init 初始化路由
func Init(r *gin.Engine) {
	// 跨域处理
	r.Use(func(c *gin.Context) {
		// 设置跨域头
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "authorization, Content-Type")

		// 处理 OPTIONS 请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// 登录相关接口
	login := r.Group("/login", midware.ErrorHandler(), midware.MongoClientMiddleware())
	{
		login.GET("/student", handler.LoginHandler_Student)
		login.GET("/counsellor", handler.LoginHandler_Counsellor)
	}

	// API 相关接口
	api := r.Group("/api", midware.ErrorHandler(), midware.MongoClientMiddleware())
	{
		// WebSocket 连接
		api.GET("/ws/:UserID", handler.WebSocketConnection)

		// 学生相关接口
		student := api.Group("/student", midware.JWTAuthMiddleware("Counsellor", "Student"), midware.GetRabbitMQMiddleware(), midware.RedisClientMiddleware())
		{
			// 修改学生个人信息
			student.PUT("/profile/:UserID", studenthandler.ModifiedProfileHandler)
			// 提交反馈或建议
			student.POST("/feedbackOradvice", studenthandler.FeedbackOAdvice)
			// 获取成绩
			student.GET("/score", studenthandler.GetConcreteSorce)
			// 提交作业
			submit := student.Group("/submit")
			{
				submit.POST("/:UserID", studenthandler.Submission)
				submit.GET("/status/:SubmissionID", studenthandler.GetSubmissionStatus)
				submit.GET("/list", studenthandler.GetSubmissionStatus)
			}
		}

		// 辅导员相关接口
		counsellor := api.Group("/counsellor", midware.JWTAuthMiddleware("Counsellor"), midware.GetRabbitMQMiddleware(), midware.RedisClientMiddleware())
		{
			// 添加事由
			counsellor.POST("/:CounsellorID/cause", counsellorhandler.AddCause)
			// 获取事由
			counsellor.GET("/:CounsellorID/cause", counsellorhandler.GetCause)
			// 设置辅导员可预约的时间
			counsellor.POST("/access-time", counsellorhandler.SetAccessTimeHandler)
			// 设置公告
			counsellor.POST("/setannouncement", counsellorhandler.SetAnnouncement)
			// 审核相关接口
			audit := counsellor.Group("/audit")
			{
				audit.GET("/list", counsellorhandler.GetSubmissionList)
				audit.PUT("/review/single", counsellorhandler.AuditOne)
				audit.PUT("/review/bulk", counsellorhandler.AuditMany)
				audit.GET("/history", counsellorhandler.GetAuditHistory)
			}
			// 学生信息相关接口
			information := counsellor.Group("/information")
			{
				information.PUT("/correct/:UserID", counsellorhandler.CorrectGrade)
				information.POST("/bulk-import/student", counsellorhandler.ImportStudent)
				information.POST("/bulk-import/counsellor", counsellorhandler.ImportCounsellor)
				information.POST("/bulk-import/mark", counsellorhandler.ImportMark)
				information.GET("/student", counsellorhandler.GetStudentInformation)
			}
		}

		// 广场相关接口
		square := api.Group("/square", midware.JWTAuthMiddleware("Counsellor", "Student"), midware.GetRabbitMQMiddleware(), midware.RedisClientMiddleware())
		{
			// 获取公告
			square.GET("/annoucement", squarehandler.GetAnnouncement)
			// 主题相关接口
			topic := square.Group("/topic")
			{
				// 发布新主题
				topic.POST("/new", squarehandler.NewTopic)
				// 修改主题
				topic.PUT("", squarehandler.ModifiedTopic)
				// 获取主题列表
				topic.GET("/list", squarehandler.GetTopicList)
				// 获取主题详情
				topic.GET("", squarehandler.GetTopic)
				// 发布回复
				topic.POST("/replies", squarehandler.NewReply)
				// 获取回复列表
				topic.GET("/replies", squarehandler.GetReply)
				// 获取浏览量和点赞量
				topic.GET("/views&likes", squarehandler.GetViewsAndlikes)
				// 点赞回复
				topic.PUT("/likes/reply", squarehandler.LikeReply)
				// 点赞主题
				topic.PUT("/likes/topic", squarehandler.LikesTopic)
				// 删除主题
				topic.DELETE("/delete/topic", squarehandler.DeleteTopic)
				// 删除回复
				topic.DELETE("/delete/reply", squarehandler.DeleteReply)
			}
		}
	}
}
