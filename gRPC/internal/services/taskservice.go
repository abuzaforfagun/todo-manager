package services

import (
	"context"
	"database/sql"
	"log"
	"time"

	task_repository "github.com/abuzaforfagun/todo-manager/internal/repositories"
	"github.com/abuzaforfagun/todo-manager/models"
	"github.com/abuzaforfagun/todo-manager/protogen/golang/task"
	"google.golang.org/genproto/googleapis/type/date"
)

type TasksService struct {
	db *sql.DB
	task.UnimplementedTasksServer
}

func NewTasksService(db *sql.DB) TasksService {
	return TasksService{db: db}
}

func (t *TasksService) AddTask(_ context.Context, req *task.TaskRequestModel) (*task.Empty, error) {
	model := models.Task{
		Name:      req.Name,
		Status:    models.Pending,
		CreatedAt: time.Now().Local().UTC(),
	}

	err := task_repository.Add(model)

	if err != nil {
		log.Printf("Unable to add %v", err)
		return &task.Empty{}, err
	}
	return &task.Empty{}, nil
}

func timeToGoogleDate(t time.Time) *date.Date {
	return &date.Date{
		Year:  int32(t.Year()),
		Month: int32(t.Month()),
		Day:   int32(t.Day()),
	}
}

func (t *TasksService) GetAll(context context.Context, req *task.Empty) (*task.TaskList, error) {
	responseFromDb, err := task_repository.GetAll(context)
	if err != nil {
		log.Printf("Unable to retrieve %v", err)
		return &task.TaskList{}, err
	}

	var tasksDto []*task.TaskResponseModel

	for _, dbTask := range responseFromDb {
		taskDto := task.TaskResponseModel{
			Id:        int32(dbTask.Id),
			Name:      dbTask.Name,
			Status:    dbTask.Status.ToString(),
			CreatedAt: timeToGoogleDate(dbTask.CreatedAt),
		}

		tasksDto = append(tasksDto, &taskDto)
	}

	result := task.TaskList{
		Tasks: tasksDto,
	}

	return &result, nil
}

func (t *TasksService) SetToInProgress(_ context.Context, req *task.IntWrapper) (*task.Empty, error) {
	err := task_repository.UpdateStatusToInProgress(int(req.Value))

	if err != nil {
		log.Printf("Error: Unable to update status %v", err)
		return &task.Empty{}, err
	}

	return &task.Empty{}, nil
}

func (t *TasksService) SetToCompleted(_ context.Context, req *task.IntWrapper) (*task.Empty, error) {
	err := task_repository.UpdateStatusToCompleted(int(req.Value))

	if err != nil {
		log.Printf("Error: Unable to update status %v", err)
		return &task.Empty{}, err
	}

	return &task.Empty{}, nil
}

func (t *TasksService) Delete(_ context.Context, req *task.IntWrapper) (*task.Empty, error) {
	err := task_repository.Delete(int(req.Value))

	if err != nil {
		log.Printf("Error: Unable to delete %v", err)
		return &task.Empty{}, err
	}

	return &task.Empty{}, nil
}
