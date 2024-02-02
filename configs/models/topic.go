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

// 数据模型记得在更新的时候要改
