package models

type SubmitInformation struct {
	CurrentUser
	ItemName     string
	AcademicYear string
	Evidence     []string
	Status       bool
}
