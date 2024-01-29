package models

type Form struct {
	formId         int64
	userId         int64
	academicYear   string
	item           Item
	status         uint8
	dismissalCause string
	advice         string
}

type Item struct {
}
