package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"todo-console/core"
	"todo-console/store_in_file"
)

type Task = core.Task

func main() {

	err := store_in_file.Init()
	if err != nil {
		fmt.Println("Unable to initialize")
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
			addTask()
		case 2:
			deleteTask()
		case 3:
			displayTaskList()
		case 4:
			defer store_in_file.CloseConnection()
			fmt.Println("----- THANK YOU -----")
			return
		default:
			fmt.Println("Please enter correct value")
		}
	}

}

func deleteTask() {
	fmt.Println("Please enter the task number: ")

	reader := bufio.NewReader(os.Stdin)
	taskToDelete, err := reader.ReadString('\n')

	taskToDelete = strings.TrimSpace(taskToDelete)

	if err != nil {
		fmt.Println("Faild to retrive task number")
	}

	if err == nil {
		err := store_in_file.DeleteTask(taskToDelete)

		if err != nil {
			fmt.Println("Failed to delete the task, ", err)
			return
		}
	}

	clearScreen()

	fmt.Println("Task deleted.")
}

func addTask() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please enter the task name")
	taskName, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Error reading input", err)
	}

	store_in_file.AddTask(taskName)

	clearScreen()

	fmt.Println("Task added.")
}

func displayMenu() {
	fmt.Println(`----- Todo Management -----
	1. Add
	2. Delete
	3. Display
	4. Exit`)
}

func displayTaskList() {
	tasks := store_in_file.GetTasks()

	fmt.Println("ID \t Name")
	for _, task := range tasks {
		task.Print()
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
