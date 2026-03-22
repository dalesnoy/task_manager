package handler

import (
	"net/http"
	"strconv"

	"task-manager/internal/model"
	"task-manager/internal/service"
	"task-manager/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// GetProjects godoc
// @Summary Список проектов пользователя
// @Tags projects
// @Security BearerAuth
// @Produce json
// @Success 200 {array} model.Project
// @Router /api/projects [get]
func GetProjects(c *gin.Context) {
	userID := middleware.GetUserID(c)
	projects, err := service.GetProjects(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, projects)
}

// GetProject godoc
// @Summary Получить проект по ID
// @Tags projects
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID проекта"
// @Success 200 {object} model.Project
// @Failure 404 {object} model.ErrorResponse
// @Router /api/projects/{id} [get]
func GetProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: "Некорректный ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	userID := middleware.GetUserID(c)
	project, err := service.GetProject(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{
			Error: err.Error(),
			Code:  http.StatusNotFound,
		})
		return
	}
	c.JSON(http.StatusOK, project)
}

// CreateProject godoc
// @Summary Создать проект
// @Tags projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body model.ProjectInput true "Данные проекта"
// @Success 201 {object} model.Project
// @Failure 400 {object} model.ErrorResponse
// @Router /api/projects [post]
func CreateProject(c *gin.Context) {
	var input model.ProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: "Некорректные данные: " + err.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	userID := middleware.GetUserID(c)
	project, err := service.CreateProject(input, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: err.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}
	c.JSON(http.StatusCreated, project)
}

// UpdateProject godoc
// @Summary Обновить проект
// @Tags projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID проекта"
// @Param input body model.ProjectInput true "Данные проекта"
// @Success 200 {object} model.Project
// @Failure 400 {object} model.ErrorResponse
// @Router /api/projects/{id} [put]
func UpdateProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: "Некорректный ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var input model.ProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: "Некорректные данные: " + err.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	userID := middleware.GetUserID(c)
	project, err := service.UpdateProject(uint(id), input, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: err.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}
	c.JSON(http.StatusOK, project)
}

// DeleteProject godoc
// @Summary Удалить проект
// @Tags projects
// @Security BearerAuth
// @Param id path int true "ID проекта"
// @Success 200 {object} map[string]string
// @Failure 400 {object} model.ErrorResponse
// @Router /api/projects/{id} [delete]
func DeleteProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: "Некорректный ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	userID := middleware.GetUserID(c)
	if err := service.DeleteProject(uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: err.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Проект удалён"})
}
