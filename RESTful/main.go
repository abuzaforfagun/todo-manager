package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	config "restful-service/configurations"
	"restful-service/db"
	auth_handlers "restful-service/handlers/auth"
	todo_handlers "restful-service/handlers/todo"
	"restful-service/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	file, err := os.Open("config.json")

	if err != nil {
		log.Fatalf("Unable to open config file %v", err)
	}
	defer file.Close()

	var config config.Configuration
	json.NewDecoder(file).Decode(&config)
	_, err = json.Marshal(config)

	if err != nil {
		log.Fatalf("Unable to decode config file %v", err)
	}
	auth_handlers.EncryptionKey = config.EncryptionKey
	auth_handlers.JwtKey = config.JwtKey

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Database.UserName, config.Database.Password, config.Database.Server, config.Database.Database)

	err = db.Init(dsn)
	if err != nil {
		log.Fatalf("Unable to open database connection %v", err)
	}
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		if c.Request.URL.Path != "/login" && c.Request.URL.Path != "/user/register" {
			middleware.AuthMiddleware()(c)
		} else {
			c.Next()
		}
	})

	router.POST("/user/register", auth_handlers.Register)
	router.POST("/login", auth_handlers.Login)
	router.GET("/todo", todo_handlers.GetAll)
	router.POST("/todo", todo_handlers.Add)
	router.POST("/todo/:id/:status", todo_handlers.UpdateStatus)
	router.DELETE("/todo/:id", todo_handlers.Delete)

	err = router.Run(":8000")

	if err != nil {
		panic(err)
	}
}
