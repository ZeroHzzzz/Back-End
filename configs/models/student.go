package models

type Student struct {
	UserId     int64  `json:"userId"`
	UserName   string `json:"userName"`
	passWord   string `json:"-"`
	Class      string `json:"class"`
	Profession string `json:"profession"`
	Grade      int    `json:"grade"`
	// FeedbackId  string `json:"feedbackId"`  //申诉
	// RecommemdId string `json:"recommemdId"` // 建议
	// FormId      string `json:"formId"`      // 申报表
}

// 最后三项放到nosql中
