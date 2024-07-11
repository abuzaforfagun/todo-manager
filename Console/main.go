package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"todo-console/core"

	// repository "todo-console/store_in_file"
	repository "todo-console/store_in_database"
)

type Task = core.Task

func main() {
	err := repository.Init()
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
			defer repository.CloseConnection()
			fmt.Println("----- THANK YOU -----")
			return
		case 5:
			updateStatus(core.InProgress)
		case 6:
			updateStatus(core.Completed)
		default:
			fmt.Println("Please enter correct value")
		}
	}

}

func deleteTask() {
	fmt.Println("Please enter the task number: ")

	reader := bufio.NewReader(os.Stdin)
	taskToDelete, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Faild to retrive task number")
		return
	}
	taskToDelete = strings.TrimSpace(taskToDelete)

	taskId, err := strconv.Atoi(taskToDelete)

	if err == nil {
		err := repository.DeleteTaskById(taskId)
		if err != nil {
			fmt.Println("Failed to delete task", err)
			return
		}
	}
	if err != nil {
		err := repository.DeleteTaskByName(taskToDelete)
		if err != nil {
			fmt.Println("Failed to delete task", err)
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

	err = repository.AddTask(taskName)

	if err != nil {
		fmt.Println("Failed to add", err)
		return
	}

	clearScreen()

	fmt.Println("Task added.")
}

func displayMenu() {
	fmt.Println(`----- Todo Management -----
	1. Add
	2. Delete
	3. Display
	4. Exit
	5. Update to In Progres
	6. Update to Completed`)
}

func displayTaskList() {
	tasks, err := repository.GetTasks()

	if err != nil {
		fmt.Println("Failed to retrive tasks", err)
		return
	}

	fmt.Println("ID \t Name")
	for _, task := range tasks {
		task.Print()
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func updateStatus(status core.TaskStatus) {
	fmt.Println("Enter task id")

	var taskId int
	reader := bufio.NewReader(os.Stdin)
	_, err := fmt.Fscanf(reader, "%d", &taskId)
	if err != nil {
		fmt.Println("Unable to get the id")
		return
	}

	var task Task

	if status == core.InProgress {
		task, err = repository.UpdateToInProgress(taskId)
	}

	if status == core.Completed {
		task, err = repository.UpdateToCompleted(taskId)
	}
	if err != nil {
		fmt.Println("Unable to update the status", err)
	}
	fmt.Printf("Task[%s] status is updated to %s\n", strings.TrimSpace(task.Name), task.Status.ToString())
}
