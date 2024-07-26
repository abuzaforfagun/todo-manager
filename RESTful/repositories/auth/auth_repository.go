package auth_repository

import (
	"errors"
	"log"
	"restful-service/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	Register(userName string, password string) error
	HasUser(username string) (bool, error)
	GetUser(username string) (user models.UserDto, err error)
}

type authRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) Register(userName string, password string) error {
	credential := models.User{
		Username: userName,
		Password: password,
	}

	result := r.db.Create(&credential)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			log.Println("ERROR: User already exists")
			return errors.New("User already exists")
		}
		return result.Error
	}

	return nil
}

func (r *authRepository) HasUser(username string) (bool, error) {
	var credential models.User
	result := r.db.Find(&credential, "Username=?", username)

	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (r *authRepository) GetUser(username string) (user models.UserDto, err error) {
	var credential models.User
	result := r.db.Find(&credential, "Username=?", username)
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
