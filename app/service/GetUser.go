package service

import (
	database "hr/configs/database/mySQL"
	"hr/configs/models"
)

func GetUser(UserId string) (models.Student, error) {
	var user models.Student
	err := database.DB.Where("userId", UserId).First(&user).Error
	return user, err
}
