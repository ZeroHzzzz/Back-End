package controller

import "github.com/gin-gonic/gin"

// type submitInformation struct {
// 	ItemName     string `json:"itemName"`
// 	AcademicYear string `json:"academicYear"`

// }

func ScoreHandler(c *gin.Context) {
	// 上传申报
	c.Header("Content-Type", "application/json")
	itemName := c.PostForm("itemName")
	academicYear := c.PostForm("academicYear")
	file, err := c.FormFile("evidence")
	if err != nil {
		// 处理逻辑
		return
	}
	currentUser, ok := c.Get("CurrentUser")
	if !ok {
		return
	}

}
