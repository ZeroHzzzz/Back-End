package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubmitInformation struct {
	CurrentUser
	ItemName     string    `bson:"itemName"`
	ItemValue    int64     `bson:"itemValue"`
	AcademicYear string    `bson:"academicYear"`
	Msg          string    `bson:"msg"`
	Evidence     []string  `bson:"evidence"`
	AduiterId    string    `bson:"aduiterId"`
	Status       bool      `bson:"status"`
	Cause        string    `bson:"cause"`
	Advice       string    `bson:"advice"`
	CreateAt     time.Time `bson:"create_at"`
}

type SubmitHistory struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	SubmissionId string             `bson:"submissionId"`
	AuditorId    string             `bson:"auditorId"` // 审核人
	Message      string             `bson:"message"`
	CreateAt     time.Time          `bson:"create_at"`
}
