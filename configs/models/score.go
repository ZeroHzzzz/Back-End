package models

type Score struct {
	UserID       string `bson:"UserID"`
	AcademicYear string `bson:"AcademicYear"`
	ItemName     string `bson:"ItemName"`
	Mark         int64  `bson:"Mark"`
	Msg          string `bson:"Msg"`
}

// 修改一下，修改成一个用户，然后有很多个文档，然后每个文档有两个数据，一个是itemName一个是grade
