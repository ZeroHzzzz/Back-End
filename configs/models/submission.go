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
	CreateAt     time.Time          `bson:"create_at"`
}