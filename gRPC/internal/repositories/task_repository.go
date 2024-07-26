package task_repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/abuzaforfagun/todo-manager/models"
)

type TaskRepository interface {
	GetAll() (tasks []models.Task, err error)
	Add(task models.Task) error
	Delete(taskId int) error
	UpdateStatusToInProgress(taskId int) error
	UpdateStatusToCompleted(taskId int) error
}

type taskRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) TaskRepository {
	return &taskRepository{
		db: db,
	}
}

func (t *taskRepository) GetAll() (tasks []models.Task, err error) {
	rows, err := t.db.Query("SELECT Id, Name, Status, CreatedAt FROM Tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (t *taskRepository) Add(task models.Task) error {
	sql, err := t.db.Prepare("INSERT INTO Tasks (Name, Status, CreatedAt) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = sql.Exec(task.Name, task.Status, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (t *taskRepository) Delete(taskId int) error {
	sql, err := t.db.Prepare("DELETE FROM Tasks WHERE Id = ?")
	if err != nil {
		return err
	}
	result, err := sql.Exec(taskId)
	if err != nil {
		return err
	}

	deletedRows, _ := result.RowsAffected()
	if deletedRows != 1 {
		return errors.New("Task not found")
	}

	return nil
}

func (t *taskRepository) UpdateStatusToInProgress(taskId int) error {
	sql, err := t.db.Prepare("UPDATE Tasks SET Status = ? WHERE Id = ?")
	if err != nil {
		return err
	}
	_, err = sql.Exec(models.InProgress, taskId)
	if err != nil {
		return err
	}

	return nil
}

func (t *taskRepository) UpdateStatusToCompleted(taskId int) error {
	sql, err := t.db.Prepare("UPDATE Tasks SET Status = ? WHERE Id = ?")
	if err != nil {
		return err
	}
	_, err = sql.Exec(models.Completed, taskId)
	if err != nil {
		return err
	}

	return nil
}
