package db

type DbConfig struct {
	Server   string `json:"server"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}
