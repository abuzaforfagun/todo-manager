package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	Id   int
	Name string
}

var taskList []Task

func main() {
	defer storeTasksInTextFile()
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
			addTask(reader)
		case 2:
			deleteTask(reader)
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

func deleteTask(reader *bufio.Reader) {
	fmt.Println("Please enter the task number: ")
	taskToDelete, err := reader.ReadString('\n')

	taskToDelete = strings.TrimRight(taskToDelete, "\n")

	if err != nil {
		fmt.Println("Error reading input", err)
	}

	taskNumber, err := strconv.Atoi(taskToDelete)

	if err == nil {
		for index, value := range taskList {
			if value.Id == taskNumber {
				taskList = append(taskList[:index], taskList[index+1:]...)
			}
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
		fmt.Println("Task not found!")
		return
	}
	taskList = append(taskList[:taskIndex], taskList[taskIndex+1:]...)
}

func addTask(reader *bufio.Reader) {
	fmt.Println("Please enter the task name")
	taskName, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Error reading input", err)
	}

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
	for _, value := range taskList {
		fmt.Printf("%d \t %s\n", value.Id, value.Name)
	}
}

func storeTasksInTextFile() {
	jsonData, err := json.Marshal(taskList)

	if err != nil {
		fmt.Println("Unable to convert to JSON")
		return
	}
	err = os.WriteFile("tasks.txt", jsonData, 0644)

	if err != nil {
		fmt.Println("Failed to store", err)
	}
}

func readExistingTasks() ([]Task, error) {
	tasksFileText, err := os.ReadFile("tasks.txt")

	if err != nil {
		file, err := os.Create("tasks.txt")
		defer file.Close()

		return []Task{}, err
	}

	err = json.Unmarshal(tasksFileText, &taskList)

	if err != nil {
		fmt.Println("Unable to deserialize")
	}
	return taskList, nil

}
