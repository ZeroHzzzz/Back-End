package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubmitInformation struct {
	CurrentUser  `bson:"currentUser"`
	ItemName     string   `bson:"ItemName"`
	ItemValue    int64    `bson:"ItemValue"`
	AcademicYear string   `bson:"AcademicYear"`
	Msg          string   `bson:"Msg,omitempty"` // 这个是提交的项目描述
	Evidence     []string `bson:"Evidence,omitempty"`
	AduiterID    string   `bson:"AduiterID,omitempty"`
	Status       bool     `bson:"Status"`
	Cause        string   `bson:"Cause,omitempty"`
	Advice       string   `bson:"Advice,omitempty"`
	CreateAt     int64    `bson:"CreateAt"`
}

type SubmitHistory struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	SubmissionID string             `bson:"SubmissionID"`
	AuditorID    string             `bson:"AuditorID"` // 审核人
	Message      string             `bson:"Message"`
	CreateAt     int64              `bson:"CreateAt"`
}
