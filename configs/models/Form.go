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

// 申请条目
type Item struct {
}
