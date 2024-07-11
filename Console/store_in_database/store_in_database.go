package store_in_database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"todo-console/core"

	_ "github.com/go-sql-driver/mysql"
)

type Task = core.Task

func Init() error { return nil }

func connectDatabase() (*sql.DB, error) {
	file, err := os.Open("config.json")

	if err != nil {
		fmt.Println("Unable to open config file", err)
	}
	defer file.Close()

	var dbConfig DbConfig
	json.NewDecoder(file).Decode(&dbConfig)
	_, err = json.Marshal(dbConfig)

	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbConfig.UserName, dbConfig.Password, dbConfig.Server, dbConfig.Database)

	// Open database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetTasks() ([]Task, error) {
	db, err := connectDatabase()

	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT Id, Name, Status FROM Tasks"

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.Id, &task.Name, &task.Status)

		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func CloseConnection() error {
	return nil
}

func AddTask(taskName string) error {
	db, err := connectDatabase()
	if err != nil {
		return err
	}
	defer db.Close()
	sql, err := db.Prepare("INSERT INTO Tasks (Name, Status, CreatedAt) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = sql.Exec(taskName, core.Pending, time.Now().UTC())

	return err
}

func UpdateToInProgress(id int) (Task, error) {
	db, err := connectDatabase()
	if err != nil {
		return Task{}, err
	}
	defer db.Close()

	sql, err := db.Prepare("UPDATE Tasks SET Status = ? WHERE Id = ?")

	_, err = sql.Exec(core.InProgress, id)

	if err != nil {
		return Task{}, err
	}
	query, err := db.Prepare("SELECT Id, Name, Status FROM Tasks WHERE Id = ?")
	if err != nil {
		return Task{}, err
	}
	queryResult := query.QueryRow(id)
	var task Task
	err = queryResult.Scan(&task.Id, &task.Name, &task.Status)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func UpdateToCompleted(id int) (Task, error) {
	db, err := connectDatabase()
	if err != nil {
		return Task{}, err
	}
	defer db.Close()

	sql, err := db.Prepare("UPDATE Tasks SET Status = ? WHERE Id = ?")
	if err != nil {
		return err
	}
	_, err = sql.Exec(core.Completed, id)

	if err != nil {
		return Task{}, err
	}
	query, err := db.Prepare("SELECT Id, Name, Status FROM Tasks WHERE Id = ?")
	if err != nil {
		return Task{}, err
	}
	queryResult := query.QueryRow(id)
	var task Task
	err = queryResult.Scan(&task.Id, &task.Name, &task.Status)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func DeleteTaskById(id int) error {
	db, err := connectDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	sql, err := db.Prepare("DELETE FROM Tasks WHERE Id = ?")

	if err != nil {
		return err
	}
	_, err = sql.Exec(id)

	return err
}

func DeleteTaskByName(name string) error {
	db, err := connectDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	sql, err := db.Prepare("DELETE Tasks WHERE Name LIKE ?")

	if err != nil {
		return err
	}
	_, err = sql.Exec(core.Completed, fmt.Sprintf("%%%s%%", name))

	return err
}
