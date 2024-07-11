package store_in_file

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"todo-console/core"
)

const fileName string = "tasks.txt"

type Task = core.Task

var Tasks []Task

func DeleteTask(taskToDelete string) error {

	taskNumber, err := strconv.Atoi(taskToDelete)
	if err == nil {
		for index, value := range Tasks {
			if value.Id == taskNumber {
				Tasks = append(Tasks[:index], Tasks[index+1:]...)
			}
		}
		return nil
	}

	taskIndex := -1
	for index, value := range Tasks {
		if strings.Contains(value.Name, taskToDelete) {
			taskIndex = index
			break
		}
	}
	if taskIndex == -1 {
		return errors.New("Task not found")
	}
	Tasks = append(Tasks[:taskIndex], Tasks[taskIndex+1:]...)

	return nil
}

func AddTask(taskName string) {
	taskName = strings.TrimRight(taskName, "\n")
	taskId := len(Tasks) + 1
	fmt.Println(taskId)
	task := Task{
		Id:   taskId,
		Name: taskName,
	}

	Tasks = append(Tasks, task)
}

func store() error {
	jsonData, err := json.Marshal(Tasks)

	if err != nil {
		return err
	}
	err = os.WriteFile(fileName, jsonData, 0644)

	if err != nil {
		return err
	}

	return err
}

func Init() error {
	tasksFileText, err := os.ReadFile(fileName)

	if err != nil {
		file, err := os.Create(fileName)
		defer file.Close()

		return err
	}

	err = json.Unmarshal(tasksFileText, &Tasks)

	if err != nil {
		return err
	}
	return nil
}

func GetTasks() []Task {
	return Tasks
}

func CloseConnection() {
	err := store()

	if err != nil {
		fmt.Println("Failed to store in the file", err)
	}
}

func UpdateToInProgress(id int) (Task, error) {
	for index, value := range Tasks {
		if value.Id == id {
			Tasks[index] = value.UpdateToInProgress()
			return value, nil
		}
	}

	return Task{}, errors.New("Task not found")
}

func UpdateToCompleted(id int) (Task, error) {
	for index, value := range Tasks {
		if value.Id == id {
			Tasks[index] = value.UpdateToCompleted()
			return value, nil
		}
	}

	return Task{}, errors.New("Task not found")
}
