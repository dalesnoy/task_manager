package handler

import (
	"net/http"
	"strconv"

	"task-manager/internal/model"
	"task-manager/internal/service"
	"task-manager/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// GetMembers возвращает участников проекта
func GetMembers(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "некорректный ID проекта"})
		return
	}

	userID := middleware.GetUserID(c)

	members, err := service.GetMembers(uint(projectID), userID)
	if err != nil {
		c.JSON(http.StatusForbidden, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

// AddMember добавляет участника в проект
func AddMember(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "некорректный ID проекта"})
		return
	}

	var input model.AddMemberInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "укажите корректный email"})
		return
	}

	userID := middleware.GetUserID(c)

	member, err := service.AddMember(uint(projectID), input, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, member)
}

// RemoveMember удаляет участника из проекта
func RemoveMember(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "некорректный ID проекта"})
		return
	}

	memberUserID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "некорректный ID пользователя"})
		return
	}

	userID := middleware.GetUserID(c)

	if err := service.RemoveMember(uint(projectID), uint(memberUserID), userID); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "участник удалён"})
}
