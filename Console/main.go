package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var taskList []string

func main() {

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
		taskIndex := taskNumber - 1
		taskList = append(taskList[:taskIndex], taskList[taskIndex+1:]...)
	}

	taskIndex := -1
	for index, value := range taskList {
		if strings.Contains(value, taskToDelete) {
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
	taskList = append(taskList, strconv.Itoa(len(taskList)+1)+". "+taskName)
}

func displayMenu() {
	fmt.Println(`----- Todo Management -----
	1. Add
	2. Delete
	3. Display
	4. Exit`)
}

func displayTaskList() {
	for _, value := range taskList {
		fmt.Printf("%s\n", value)
	}
}
