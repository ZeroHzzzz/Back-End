package models

type Student struct {
	UserId      int64          `bson:"userId"`
	UserName    string         `bson:"userName"`
	passWord    string         `bson:"-"`
	Class       string         `bson:"class"`
	Profession  string         `bson:"profession"`
	Grade       int            `bson:"grade"` // 这里是年级
	Mark        map[string]int `bson:"mark"`
	FeedbackId  []string       `bson:"feedbackId"`  //申诉
	RecommemdId []string       `bson:"recommemdId"` // 建议
	FormId      []string       `bson:"formId"`      // 申报表
}

type CurrentUser struct {
	UserId     string
	UserName   string
	Grade      int //这里是年级不是成绩
	Role       string
	Profession string
}
