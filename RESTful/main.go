package main

import (
	"encoding/json"
	"fmt"
	"os"
	"restful-service/db"
	todo_handlers "restful-service/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	file, err := os.Open("config.json")

	if err != nil {
		fmt.Println("Unable to open config file", err)
	}
	defer file.Close()

	var dbConfig db.DbConfig
	json.NewDecoder(file).Decode(&dbConfig)
	_, err = json.Marshal(dbConfig)

	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbConfig.UserName, dbConfig.Password, dbConfig.Server, dbConfig.Database)

	err = db.Init(dsn)
	if err != nil {
		panic(err)
	}
	router := gin.Default()

	router.GET("/todo", todo_handlers.GetAll)
	router.POST("/todo", todo_handlers.Add)
	router.POST("/todo/:id/:status", todo_handlers.UpdateStatus)
	router.DELETE("/todo/:id", todo_handlers.Delete)

	err = router.Run(":8000")

	if err != nil {
		panic(err)
	}
}
