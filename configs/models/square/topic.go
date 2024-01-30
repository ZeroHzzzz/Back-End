package square

import "time"

type Topic struct {
	TopicId   string    `bson:"_id,omitempty"`
	Title     string    `bson:"title"`
	Content   string    `bson:"content"`
	AutherId  string    `bson:"autherId"`
	ReplyId   string    `bson:"replyId"`
	Likes     int       `bson:"likes"`    //点赞量
	ViewTime  int       `bson:"viewTime"` //浏览量
	CreatedAt time.Time `bson:"createdAt"`
}

type Reply struct {
	ItemId    string    `bson:"itemId"`  //回复本身的id
	ReplyId   string    `bson:"replyId"` // 回复的上一级的id
	Content   string    `bson:"content"`
	ReplyerId string    `bson:"replyerId"`
	Likes     int       `bson:"likes"` //点赞量
	CreatedAt time.Time `bson:"createdAt"`
}
