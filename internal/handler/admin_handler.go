package handler

import (
	"net/http"
	"strconv"

	"task-manager/internal/model"
	"task-manager/internal/repository"

	"github.com/gin-gonic/gin"
)

// AdminGetUsers возвращает список всех пользователей
func AdminGetUsers(c *gin.Context) {
	users, err := repository.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "ошибка получения пользователей"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// AdminGetUserProjects возвращает все проекты конкретного пользователя (свои + участие)
func AdminGetUserProjects(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "некорректный ID пользователя"})
		return
	}

	projects, err := repository.GetProjectsByUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "ошибка получения проектов"})
		return
	}

	c.JSON(http.StatusOK, projects)
}
