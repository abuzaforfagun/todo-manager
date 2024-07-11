package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const fileName string = "tasks.txt"

var taskList []Task

func main() {
	defer func() {
		err := storeTasksInTextFile()

		if err != nil {
			fmt.Println("Failed to store in the file", err)
		}
	}()
	var readExistingTaskErr error
	taskList, readExistingTaskErr = readExistingTasks()

	if readExistingTaskErr != nil {
		fmt.Println("Unable to retrieve tasks", readExistingTaskErr)
	}

	for {
		displayMenu()
		var input int32

		reader := bufio.NewReader(os.Stdin)
		_, err := fmt.Fscanf(reader, "%d\n", &input)

		if err != nil {
			fmt.Println("Error reading input", err)
		}

		switch input {
		case 1:
			fmt.Println("Please enter the task name")
			taskName, err := reader.ReadString('\n')

			if err != nil {
				fmt.Println("Error reading input", err)
			}
			addTask(taskName)
		case 2:
			fmt.Println("Please enter the task number: ")
			taskToDelete, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Faild to retrive task number")
			}
			if err == nil {
				err := deleteTask(taskToDelete)

				if err != nil {
					fmt.Println("Failed to delete the task, ", err)
				}
			}

		case 3:
			displayTaskList()
		case 4:
			fmt.Println("----- THANK YOU -----")
			return
		default:
			fmt.Println("Please enter correct value")
		}
	}
}

func deleteTask(taskToDelete string) error {
	taskToDelete = strings.TrimRight(taskToDelete, "\n")

	taskNumber, err := strconv.Atoi(taskToDelete)

	if err != nil {
		return err
	}

	for index, value := range taskList {
		if value.Id == taskNumber {
			taskList = append(taskList[:index], taskList[index+1:]...)
		}
	}

	taskIndex := -1
	for index, value := range taskList {
		if strings.Contains(value.Name, taskToDelete) {
			taskIndex = index
			break
		}
	}
	if taskIndex == -1 {
		return errors.New("Task not found")
	}
	taskList = append(taskList[:taskIndex], taskList[taskIndex+1:]...)

	return nil
}

func addTask(taskName string) {
	taskName = strings.TrimRight(taskName, "\n")
	taskId := len(taskList) + 1
	task := Task{
		Id:   taskId,
		Name: taskName,
	}

	taskList = append(taskList, task)
}

func displayMenu() {
	fmt.Println(`----- Todo Management -----
	1. Add
	2. Delete
	3. Display
	4. Exit`)
}

func displayTaskList() {
	fmt.Println("ID \t Name")
	for _, task := range taskList {
		task.Print()
	}
}

func storeTasksInTextFile() error {
	jsonData, err := json.Marshal(taskList)

	if err != nil {
		return err
	}
	err = os.WriteFile(fileName, jsonData, 0644)

	if err != nil {
		return err
	}

	return err
}

func readExistingTasks() ([]Task, error) {
	tasksFileText, err := os.ReadFile(fileName)

	if err != nil {
		file, err := os.Create(fileName)
		defer file.Close()

		return nil, err
	}

	err = json.Unmarshal(tasksFileText, &taskList)

	if err != nil {
		return nil, err
	}
	return taskList, nil
}
