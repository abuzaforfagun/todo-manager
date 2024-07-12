package db

import (
	"database/sql"
	"fmt"

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

	fmt.Println("Database connection established successfully!")
	return nil
}

func Get() *sql.DB {
	return db
}
func Close() {
	db.Close()
	fmt.Println("Database connection closed.")
}
