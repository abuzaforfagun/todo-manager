package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/abuzaforfagun/todo-manager/internal/db"
	task_repository "github.com/abuzaforfagun/todo-manager/internal/repositories"
	"github.com/abuzaforfagun/todo-manager/internal/services"
	"github.com/abuzaforfagun/todo-manager/protogen/golang/task"
	"google.golang.org/grpc"
)

func main() {

	file, err := os.Open("../../config.json")

	if err != nil {
		log.Fatalf("Unable to open config file %v", err)
	}
	defer file.Close()

	var config db.DbConfig
	json.NewDecoder(file).Decode(&config)
	_, err = json.Marshal(config)

	if err != nil {
		log.Fatalf("Unable to decode config file %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.UserName, config.Password, config.Server, config.Database)

	err = db.Init(dsn)

	db := db.Get()

	const addr = "localhost:50051"

	// create a TCP listener on the specified port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	taskRepository := task_repository.NewRepository(db)
	taskService := services.NewTasksService(taskRepository)

	task.RegisterTasksServer(server, &taskService)

	log.Printf("server listening at %v", listener.Addr())
	if err = server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
