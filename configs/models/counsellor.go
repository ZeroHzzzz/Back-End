package models

type Counsellor struct {
	UserID     string `bson:"_id"`
	UserName   string `bson:"UserName"`
	PassWord   string `bson:"PassWord"`
	Grade      string `bson:"Grade"`
	Profession string `bson:"Profession"`
}
