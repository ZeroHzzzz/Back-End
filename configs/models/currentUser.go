package models

type CurrentUser struct {
	UserID     string `bson:"UserID"`
	UserName   string `bson:"UserName"`
	Grade      string `bson:"Grade"` //这里是年级不是成绩
	Role       string `bson:"Role"`
	Profession string `bson:"Profession"`
}
