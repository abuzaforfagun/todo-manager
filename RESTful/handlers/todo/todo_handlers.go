package todo_handlers

import (
	"context"
	"log"
	"net/http"
	"restful-service/models"
	todo_repositories "restful-service/repositories/todo"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Get Todo
// @Description Get todo list
// @Tags todo
// @Produce json
// @Security BearerAuth
// @Param pageSize query int false "Page size (Default 10)"
// @Param pageNumber query int false "Page number (Default 1)"
// @Success 200 {object} []models.TaskDto
// @Router /todo [get]
func GetAll(ctx context.Context, c *gin.Context) {
	pageSizeParam, hasPageSize := c.GetQuery("pageSize")

	pageNumberParam, hasPageNumber := c.GetQuery("pageNumber")

	var pageSize = 10
	var pageNumber = 1

	if hasPageSize {
		var err error
		pageSize, err = strconv.Atoi(pageSizeParam)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please verify pageSize"})
		}
	}

	if hasPageNumber {
		var err error
		pageNumber, err = strconv.Atoi(pageNumberParam)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please verify pageNumber"})
		}
	}

	userId := c.GetUint("UserId")
	result, err := todo_repositories.GetAll(userId, pageSize, pageNumber)

	if err != nil {
		log.Printf("Error: Unable to get todo list %v", err)
		c.JSON(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, result)
}

// @Summary Add Todo
// @Description Add new todo item
// @Tags todo
// @Produce json
// @Accept json
// @Security BearerAuth
// @Param todo body models.TaskRequestDto true "Task payload"
// @Success 201
// @Router /todo [post]
func Add(c *gin.Context) {
	var task models.TaskRequestDto

	userId := c.GetUint("UserId")
	err := c.BindJSON(&task)
	if err != nil {
		log.Printf("Warning: Invalid request %v", err)
		c.JSON(http.StatusBadRequest, nil)
	}

	err = todo_repositories.Add(task, userId)
	if err != nil {
		log.Printf("Error: Unable to add todo %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusAccepted, gin.H{})
}

// @Summary Delete Todo
// @Description Delete todo item
// @Tags todo
// @Produce json
// @Security BearerAuth
// @Param todo query int true "Todo id to delete"
// @Success 200
// @Router /todo [delete]
func Delete(c *gin.Context) {
	idParam := c.Query("id")

	userId := c.GetUint("UserId")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		log.Printf("Warning: Invalid request %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = todo_repositories.Delete(id, userId)

	if err != nil {
		log.Printf("Error: Unable to delete %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func UpdateStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	userId := c.GetUint("UserId")

	if err != nil {
		log.Printf("Warning: Invalid request %v", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	status := c.Param("status")

	if status == "inprogress" {
		err = todo_repositories.UpdateStatusToInProgress(id, userId)
	}
	if status == "completed" {
		err = todo_repositories.UpdateStatusToCompleted(id, userId)
	}

	if err != nil {
		log.Printf("Error: Unable to update status %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{})
}
