package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubmitInformation struct {
	CurrentUser  `bson:"currentUser"`
	ItemName     string    `bson:"itemName"`
	ItemValue    int64     `bson:"itemValue"`
	AcademicYear string    `bson:"academicYear"`
	Msg          string    `bson:"msg,omitempty"` // 这个是提交的项目描述
	Evidence     []string  `bson:"evidence,omitempty"`
	AduiterId    string    `bson:"aduiterId,omitempty"`
	Status       bool      `bson:"status"`
	Cause        string    `bson:"cause,omitempty"`
	Advice       string    `bson:"advice,omitempty"`
	CreateAt     time.Time `bson:"create_at"`
}

type SubmitHistory struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	SubmissionId string             `bson:"submissionId"`
	AuditorId    string             `bson:"auditorId"` // 审核人
	Message      string             `bson:"message"`
	CreateAt     time.Time          `bson:"create_at"`
}
