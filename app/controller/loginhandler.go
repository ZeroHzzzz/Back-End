package controller

import "github.com/gin-gonic/gin"

func loginHandler(c *gin.Context) {
	userid := c.PostForm("userId")
	password := c.PostForm("passWord")

	c.Set("currentUser", user) //将用户信息储存到上下文
}
