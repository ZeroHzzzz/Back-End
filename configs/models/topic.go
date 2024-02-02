package models

import "time"

type Topic struct {
	TopicId   string    `bson:"_id,omitempty"`
	Title     string    `bson:"title"`
	Content   string    `bson:"content"`
	AutherId  string    `bson:"autherId"`
	Likes     int       `bson:"likes"`    //点赞量
	ViewTime  int       `bson:"viewTime"` //浏览量
	CreatedAt time.Time `bson:"createdAt"`
}

type Reply struct {
	ReplyId   string    `bson:"_id"` //回复本身的id
	TopicId   string    `bson:"topicId"`
	ParentId  string    `bson:"parentId,omitempty"` // 回复的上一级的id
	SonId     bool      `bson:"sonId,omitempty"`    // 有子回复就更新为true
	Content   string    `bson:"content"`
	UserId    string    `bson:"userId"`
	Likes     int       `bson:"likes"` //点赞量
	CreatedAt time.Time `bson:"createdAt"`
}

// 数据模型记得在更新的时候要改
