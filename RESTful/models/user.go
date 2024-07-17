package models

import (
	"gorm.io/gorm"
)

type UserDto struct {
	UserId   uint   `json:"UserId"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserLoginDto represents a user login object
// @Description UserLoginDto represents a user login object
type UserLoginDto struct {
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

// LoginResponse represents a user login response object
// @Description LoginResponse represents a user login response object
type LoginResponse struct {
	Token string
}
