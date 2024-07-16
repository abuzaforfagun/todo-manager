package models

import (
	"gorm.io/gorm"
)

type CredentialDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Credential struct {
	gorm.Model
	Username string
	Password string
}

func (Credential) TableName() string {
	return "Credentials"
}
