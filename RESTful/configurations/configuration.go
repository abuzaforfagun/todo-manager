package config

import "restful-service/db"

type Configuration struct {
	Database      db.DbConfig `json:"database"`
	EncryptionKey string      `json:"encryptionKey"`
	JwtKey        string      `json:jwtKey`
}

func (c Configuration) GetEncryptionKey() string {
	return c.EncryptionKey
}
