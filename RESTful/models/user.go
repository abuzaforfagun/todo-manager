package models

import (
	"gorm.io/gorm"
)

type UserDto struct {
	UserId   uint   `json:"UserId"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	gorm.Model
	Username string
	Password string
	Tasks    []Task `gorm:"foreignkey:UserId"`
}

func (User) TableName() string {
	return "Users"
}
