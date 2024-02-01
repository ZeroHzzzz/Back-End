package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubmitInformation struct {
	CurrentUser
	ItemName     string   `bson:"itemName"`
	AcademicYear string   `bson:"academicYear"`
	Evidence     []string `bson:"evidence"`
	Status       bool     `bson:"status"`
}

type SubmitHistory struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	SubmissionId string             `bson:"submissionId"`
	AuditorId    string             `bson:"auditorId"`
	Status       bool               `bson:"status"`
	Cause        string             `bson:"cause"`
	Advice       string             `bson:"advice"`
	CreateAt     time.Time          `bson:"create_at"`
}

// 考虑一下要不要合并，还有history当中状态感觉应该得有三个，接受，驳回，撤回
