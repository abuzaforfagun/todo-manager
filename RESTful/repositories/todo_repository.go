package todo_repositories

import (
	"restful-service/db"
	"restful-service/models"
	"time"
)

func GetAll() ([]models.Task, error) {
	dbConnection := db.Get()
	rows, err := dbConnection.Query("SELECT Id, Name, Status, CreatedAt FROM Tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var createdAtStr string
		err := rows.Scan(&task.Id, &task.Name, &task.Status, &createdAtStr)

		if err != nil {
			return nil, err
		}
		const MySQLDateTimeLayout = "2006-01-02 15:04:05"

		task.CreatedAt, err = time.Parse(MySQLDateTimeLayout, createdAtStr)
		if err != nil {
			panic(err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func Add(task models.Task) error {
	dbConnection := db.Get()
	sql, err := dbConnection.Prepare("INSERT INTO Tasks (Name, Status, CreatedAt) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = sql.Exec(task.Name, task.Status, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func Delete(taskId int) error {
	dbConnection := db.Get()
	sql, err := dbConnection.Prepare("DELETE FROM Tasks WHERE Id = ?")
	if err != nil {
		return err
	}
	_, err = sql.Exec(taskId)
	if err != nil {
		return err
	}

	return nil
}

func UpdateStatusToInProgress(taskId int) error {
	dbConnection := db.Get()
	sql, err := dbConnection.Prepare("UPDATE Tasks SET Status = ? WHERE Id = ?")
	if err != nil {
		return err
	}
	_, err = sql.Exec(models.InProgress, taskId)
	if err != nil {
		return err
	}

	return nil
}

func UpdateStatusToCompleted(taskId int) error {
	dbConnection := db.Get()
	sql, err := dbConnection.Prepare("UPDATE Tasks SET Status = ? WHERE Id = ?")
	if err != nil {
		return err
	}
	_, err = sql.Exec(models.Completed, taskId)
	if err != nil {
		return err
	}

	return nil
}
