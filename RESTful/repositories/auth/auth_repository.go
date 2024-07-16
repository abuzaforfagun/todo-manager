package auth_repository

import (
	"errors"
	"log"
	"restful-service/db"
	"restful-service/models"

	"gorm.io/gorm"
)

func Register(userName string, password string) error {
	gormDb := db.GetGormDb()

	credential := models.User{
		Username: userName,
		Password: password,
	}

	result := gormDb.Create(&credential)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			log.Println("ERROR: User already exists")
			return errors.New("User already exists")
		}
		return result.Error
	}

	return nil
}

func HasUser(username string) (bool, error) {
	gormDb := db.GetGormDb()

	var credential models.User
	result := gormDb.Find(&credential, "Username=?", username)

	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func GetUser(username string) (user models.UserDto, err error) {
	gormDb := db.GetGormDb()
	var credential models.User
	result := gormDb.Find(&credential, "Username=?", username)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.UserDto{}, errors.New("user not found")
		}
		return models.UserDto{}, result.Error
	}

	userDto := models.UserDto{
		UserId:   credential.ID,
		Username: credential.Username,
		Password: credential.Password,
	}

	return userDto, nil
}
