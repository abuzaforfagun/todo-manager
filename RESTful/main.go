package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	config "restful-service/configurations"
	"restful-service/db"
	auth_handlers "restful-service/handlers/auth"
	todo_handlers "restful-service/handlers/todo"
	"restful-service/middleware"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
		if c.Request.URL.Path != "/login" && !strings.Contains(c.Request.URL.Path, "/swagger/") && c.Request.URL.Path != "/user/register" {
			middleware.AuthMiddleware()(c)
		} else {
			c.Next()
		}
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/user/register", auth_handlers.Register)
	router.POST("/login", auth_handlers.Login)
	router.GET("/todo", wrapHandlerWithContext(todo_handlers.GetAll))
	router.POST("/todo", todo_handlers.Add)
	router.POST("/todo/:id/:status", todo_handlers.UpdateStatus)
	router.DELETE("/todo/:id", todo_handlers.Delete)

	err = router.Run(":8000")

	if err != nil {
		panic(err)
	}
}

func wrapHandlerWithContext(handler func(ctx context.Context, c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with timeout or from request context
		ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
		defer cancel()

		// Call the handler with the context
		handler(ctx, c)
	}
}
