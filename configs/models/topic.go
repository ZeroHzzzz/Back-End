package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Topic struct {
	TopicID   primitive.ObjectID `bson:"_id"`
	Title     string             `bson:"title"`
	Content   string             `bson:"content"`
	AutherID  string             `bson:"autherId"`
	Likes     int                `bson:"likes"` //点赞量
	Views     int                `bson:"views"` //浏览量
	CreatedAt time.Time          `bson:"created_at"`
}

// 数据模型记得在更新的时候要改
