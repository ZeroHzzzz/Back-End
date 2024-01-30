package models

// 用item来记录分数
type Items struct {
	ItemName string `bson:"itemName"`
	Tag      string `bson:"tag"`
	Score    int    `bson:"score"`
}
