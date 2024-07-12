package todo_handlers

import (
	"net/http"
	"restful-service/models"

	"github.com/gin-gonic/gin"
)

func GetAll(c *gin.Context) {
	c.JSON(http.StatusOK, []models.Task{})
}
