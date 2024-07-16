package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	UserId   uint
	Username string `json:"username"`
	jwt.StandardClaims
}
