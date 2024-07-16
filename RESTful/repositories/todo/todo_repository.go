package todo_repositories

import (
	"context"
	"errors"
	"restful-service/db"
	"restful-service/models"

	"gorm.io/gorm"
)

func GetAll(ctx context.Context) ([]models.TaskDto, error) {
	gormDb := db.GetGormDb()

	var tasks []models.Task

	result := gormDb.Find(&tasks)

	if result.Error != nil {
		return nil, result.Error
	}

	var tasksDto []models.TaskDto

	for _, task := range tasks {
		taskDto := models.TaskDto{
			Id:        task.ID,
			Name:      task.Name,
			Status:    task.Status.ToString(),
			CreatedAt: task.CreatedAt,
		}

		tasksDto = append(tasksDto, taskDto)
	}

	return tasksDto, nil
}

func Add(task models.TaskDto) error {
	gormDb := db.GetGormDb()

	model := models.Task{
		Name: task.Name,
	}

	result := gormDb.Create(&model)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func Delete(taskId int) error {
	gormDb := db.GetGormDb()
	var task models.Task

	result := gormDb.Delete(&task).Where("Id=?", taskId)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateStatusToInProgress(taskId int) error {
	gormDb := db.GetGormDb()
	var task models.Task

	result := gormDb.First(&task, "Id=?", taskId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("invalid task id")
		}
		return result.Error
	}

	task.Status = models.InProgress

	return nil
}

func UpdateStatusToCompleted(taskId int) error {
	gormDb := db.GetGormDb()
	var task models.Task

	result := gormDb.First(&task, "Id=?", taskId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("invalid task id")
		}
		return result.Error
	}

	task.Status = models.Completed

	return nil
}
