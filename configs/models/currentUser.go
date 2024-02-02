package models

type CurrentUser struct {
	UserId     string
	UserName   string
	Grade      int //这里是年级不是成绩
	Role       string
	Profession string
}
