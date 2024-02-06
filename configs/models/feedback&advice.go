package models

type Feedback struct {
	feedbackId   string `bson:"_id"`
	Category     string `bson:"category"` //Advice or Feedback
	UserId       string `bson:"userId"`
	Content      string `bson:"content"`
	Status       bool   `bson:"status"`
	ReplyMessage string `bson:"replyMessage,omitempty"`
}
