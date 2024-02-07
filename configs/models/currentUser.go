package models

type CurrentUser struct {
	UserId     string
	UserName   string
	Grade      string //这里是年级不是成绩
	Role       string
	Profession string
}
