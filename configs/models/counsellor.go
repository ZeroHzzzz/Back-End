package models

type Counsellor struct {
	UserId     string `bson:"_id"`
	UserName   string `bson:"userName"`
	passWord   string `bson:"-"`
	Grade      string `bson:"grade"`
	Profession string `bson:"profession"`
}
