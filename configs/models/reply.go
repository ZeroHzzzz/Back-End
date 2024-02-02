package models

import "time"

type Reply struct {
	ReplyId   string    `bson:"_id"` //回复本身的id
	TopicId   string    `bson:"topicId"`
	ParentId  string    `bson:"parentId,omitempty"` // 回复的上一级的id
	Content   string    `bson:"content"`
	AutherId  string    `bson:"autherId"`
	Likes     int       `bson:"likes"` //点赞量
	CreatedAt time.Time `bson:"createdAt"`
}
