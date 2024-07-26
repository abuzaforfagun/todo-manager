package todo_repositories

import (
	"errors"
	"restful-service/db"
	"restful-service/models"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetAll(userId uint, pageSize int, pageNumber int) ([]models.TaskDto, error)
	Add(task models.TaskRequestDto, userId uint) error
	Delete(taskId int, userId uint) error
	UpdateStatusToInProgress(taskId int, userId uint) error
	UpdateStatusToCompleted(taskId int, userId uint) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{
		db: db,
	}
}

func (t *taskRepository) GetAll(userId uint, pageSize int, pageNumber int) ([]models.TaskDto, error) {
	gormDb := db.GetGormDb()

	var tasks []models.Task

	result := gormDb.Debug().Where(&models.Task{UserId: uint(userId)}, userId).Offset(pageSize * (pageNumber - 1)).Limit(pageSize).Find(&tasks)

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

func (t *taskRepository) Add(task models.TaskRequestDto, userId uint) error {
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

func (t *taskRepository) Delete(taskId int, userId uint) error {
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

func (t *taskRepository) UpdateStatusToInProgress(taskId int, userId uint) error {
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

func (t *taskRepository) UpdateStatusToCompleted(taskId int, userId uint) error {
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
