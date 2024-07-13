package todo_handlers

import (
	"log"
	"net/http"
	"restful-service/models"
	todo_repositories "restful-service/repositories/todo"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAll(c *gin.Context) {
	result, err := todo_repositories.GetAll()

	if err != nil {
		log.Printf("Error: Unable to get todo list %v", err)
		c.JSON(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, result)
}

func Add(c *gin.Context) {
	var task models.Task

	err := c.BindJSON(&task)
	if err != nil {
		log.Printf("Warning: Invalid request %v", err)
		c.JSON(http.StatusBadRequest, nil)
	}

	err = todo_repositories.Add(task)
	if err != nil {
		log.Printf("Error: Unable to add todo %v", err)
		c.JSON(http.StatusBadRequest, nil)
	}

	c.JSON(http.StatusAccepted, gin.H{})
}

func Delete(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		log.Printf("Warning: Invalid request %v", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	err = todo_repositories.Delete(id)

	if err != nil {
		log.Printf("Error: Unable to delete %v", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func UpdateStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		log.Printf("Warning: Invalid request %v", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	status := c.Param("status")

	if status == "inprogress" {
		err = todo_repositories.UpdateStatusToInProgress(id)
	}
	if status == "completed" {
		err = todo_repositories.UpdateStatusToCompleted(id)
	}

	if err != nil {
		log.Printf("Error: Unable to update status %v", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{})
}
