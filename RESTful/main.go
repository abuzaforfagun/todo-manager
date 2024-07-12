package main

import (
	todo_handlers "restful-service/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/todo", todo_handlers.GetAll)

	err := router.Run(":8000")

	if err != nil {
		panic(err)
	}
}
