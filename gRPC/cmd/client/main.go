package main

import (
	"context"
	"log"

	"github.com/abuzaforfagun/todo-manager/protogen/golang/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	grpcConnection, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Unable to create client %v", err)
	}

	client := task.NewTasksClient(grpcConnection)

	taskPayload := task.TaskRequestModel{
		Name: "Task number 1",
	}

	_, err = client.AddTask(context.Background(), &taskPayload)

	taskListResponse, _ := client.GetAll(context.Background(), &task.Empty{})
	taskList := taskListResponse.GetTasks()
	log.Printf("List of tasks: %v", &taskList)

	taskId := task.IntWrapper{
		Value: taskList[0].Id,
	}
	_, err = client.SetToInProgress(context.Background(), &taskId)
}
