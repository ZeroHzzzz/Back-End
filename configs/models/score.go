package models

type Score struct {
	UserId       string `bson:"userId"`
	AcademicYear string `bson:"academicYear"`
	ItemName     string `bson:"itemName"`
	Mark         int64  `bson:"mark"`
	Msg          string `bson:"msg"`
}

// 修改一下，修改成一个用户，然后有很多个文档，然后每个文档有两个数据，一个是itemName一个是grade
