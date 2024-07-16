package todo_repositories

import (
	"errors"
	"restful-service/db"
	"restful-service/models"

	"gorm.io/gorm"
)

func GetAll(userId int) ([]models.TaskDto, error) {
	gormDb := db.GetGormDb()

	var tasks []models.Task

	result := gormDb.Where(&models.Task{UserId: uint(userId)}, userId).Find(&tasks)

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

func Add(task models.TaskRequestDto, userId uint) error {
	gormDb := db.GetGormDb()

	model := models.Task{
		Name:   task.Name,
		UserId: userId,
	}

	result := gormDb.Create(&model)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func Delete(taskId int, userId uint) error {
	gormDb := db.GetGormDb()
	var task models.Task

	result := gormDb.Where("Id=?", taskId).Where("user_id=?", userId).Delete(&task)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("failed to delete")
	}

	return nil
}

func UpdateStatusToInProgress(taskId int, userId uint) error {
	gormDb := db.GetGormDb()
	var task models.Task

	result := gormDb.Where("Id=? AND user_id", taskId, userId).First(&task)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("invalid task id")
		}
		return result.Error
	}

	task.Status = models.InProgress

	return nil
}

func UpdateStatusToCompleted(taskId int, userId uint) error {
	gormDb := db.GetGormDb()
	var task models.Task

	result := gormDb.Where("Id=? AND user_id", taskId, userId).First(&task)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("invalid task id")
		}
		return result.Error
	}

	task.Status = models.Completed

	return nil
}
