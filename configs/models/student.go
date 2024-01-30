package models

type Student struct {
	UserId      int64  `json:"-"`
	UserName    string `json:"-"`
	passWord    string `json:"-"`
	Class       string `json:"-"`
	Profession  string `json:"-"`
	Grade       int    `json:"-"`
	Mark        []int  `json:"-"` // 分数。可有可无
	FeedbackId  string `json:"-"` //申诉
	RecommemdId string `json:"-"` // 建议
	FormId      string `json:"-"` // 申报表
}
