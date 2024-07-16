package models

import (
	"gorm.io/gorm"
)

type UserDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	gorm.Model
	Username string
	Password string
}

func (User) TableName() string {
	return "Users"
}
