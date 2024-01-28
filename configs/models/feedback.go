package models

type Feedback struct {
	feedbackId   int64
	userId       int64
	status       uint8
	replyMessage string
}
