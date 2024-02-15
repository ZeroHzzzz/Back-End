package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reply struct {
	ReplyID  primitive.ObjectID `bson:"_id"` //回复本身的id
	TopicID  string             `bson:"TopicID"`
	ParentID string             `bson:"ParentID,omitempty"` // 回复的上一级的id
	Content  string             `bson:"Content"`
	AutherID string             `bson:"AutherID"`
	Likes    int                `bson:"Likes"` //点赞量
	CreateAt int64              `bson:"CreateAt"`
}
