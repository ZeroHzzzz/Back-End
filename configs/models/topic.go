package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Topic struct {
	TopicID  primitive.ObjectID `bson:"_id"`
	Title    string             `bson:"Title"`
	Content  string             `bson:"Content"`
	AutherID string             `bson:"AutherID"`
	Likes    int                `bson:"Likes"` //点赞量
	Views    int                `bson:"Views"` //浏览量
	CreateAt int64              `bson:"CreateAt"`
}

// 数据模型记得在更新的时候要改
