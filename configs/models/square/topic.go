package square

import "time"

type Topic struct {
	TopicId   string    `json:"topic"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AutherId  string    `json:"autherId"`
	ReplyId   string    `json:"replyId"`
	Likes     int       `json:"likes"`    //点赞量
	ViewTime  int       `json:"viewTime"` //浏览量
	CreatedAt time.Time `json:"createdAt"`
}

type Reply struct {
	ItemId    string    `json:"itemId"`  //回复本身的id
	ReplyId   string    `json:"replyId"` // 回复的上一级的id
	Content   string    `json:"content"`
	ReplyerId string    `json:"replyerId"`
	Likes     int       `json:"likes"` //点赞量
	CreatedAt time.Time `json:"createdAt"`
}
