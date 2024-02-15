package models

type Feedback struct {
	FeedbackID   string `bson:"_id"`
	Category     string `bson:"Category"` //Advice or Feedback
	UserID       string `bson:"UserID"`
	Content      string `bson:"Content"`
	Status       bool   `bson:"Status"`
	ReplyMessage string `bson:"ReplyMessage,omitempty"`
}
