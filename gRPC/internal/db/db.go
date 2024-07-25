package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Init(dataSourceName string) error {
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return fmt.Errorf("error opening database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return fmt.Errorf("error connecting to database: %v", err)
	}

	log.Println("Database connection established")
	return nil
}

func Get() *sql.DB {
	return db
}
func Close() {
	db.Close()
	log.Println("Database connection closed")
}
