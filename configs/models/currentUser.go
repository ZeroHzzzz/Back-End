package models

type CurrentUser struct {
	UserId     string `bson:"userId"`
	UserName   string `bson:"userName"`
	Grade      string `bson:"grade"` //这里是年级不是成绩
	Role       string `bson:"role"`
	Profession string `bson:"profession"`
}
