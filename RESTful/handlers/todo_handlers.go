package todo_handlers

import (
	"net/http"
	"restful-service/models"
	todo_repositories "restful-service/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAll(c *gin.Context) {
	result, err := todo_repositories.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func Add(c *gin.Context) {
	var task models.Task

	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	err = todo_repositories.Add(task)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{})
}

func Delete(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}

	err = todo_repositories.Delete(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func UpdateStatus(c *gin.Context) {
	var task models.Task

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}
	status := c.Param("status")

	if status == "inprogress" {
		todo_repositories.UpdateStatusToInProgress(id)
	}
	if status == "completed" {
		todo_repositories.UpdateStatusToCompleted(id)
	}

	err = c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	err = todo_repositories.Add(task)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{})
}
