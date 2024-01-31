package models

type Student struct {
	UserId      int64    `bson:"userId"`
	UserName    string   `bson:"userName"`
	passWord    string   `bson:"-"`
	Class       string   `bson:"class"`
	Profession  string   `bson:"profession"`
	Grade       int      `bson:"grade"`
	FeedbackId  []string `bson:"feedbackId"`  //申诉
	RecommemdId []string `bson:"recommemdId"` // 建议
	FormId      []string `bson:"formId"`      // 申报表
}
