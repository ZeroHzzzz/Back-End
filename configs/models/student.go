package models

type Student struct {
	userId      int64  `json:"-"`
	userName    string `json:"-"`
	passWord    string `json:"-"`
	class       string `json:"-"`
	profession  string `json:"-"`
	grade       int    `json:"-"`
	mark        []int  `json:"-"` // 分数。可有可无
	feedbackId  string `json:"-"` //申诉
	recommemdId string `json:"-"` // 建议
	formId      string `json:"-"` // 申报表
}
