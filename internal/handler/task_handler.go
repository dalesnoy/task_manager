package handler

import (
	"net/http"
	"strconv"

	"task-manager/internal/model"
	"task-manager/internal/service"
	"task-manager/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// GetTasks godoc
// @Summary Список задач проекта
// @Tags tasks
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID проекта"
// @Param status query string false "Фильтр по статусу (todo, in_progress, done)"
// @Param priority query string false "Фильтр по приоритету (low, medium, high)"
// @Param sort query string false "Сортировка (deadline, priority, created_at)"
// @Param page query int false "Страница"
// @Param limit query int false "Количество на странице"
// @Success 200 {object} map[string]interface{}
// @Router /api/projects/{id}/tasks [get]
func GetTasks(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Некорректный ID", Code: 400})
		return
	}

	var filter model.TaskFilter
	c.ShouldBindQuery(&filter)

	userID := middleware.GetUserID(c)
	tasks, total, err := service.GetTasks(uint(projectID), userID, filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error(), Code: 400})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"total": total,
		"page":  filter.Page,
		"limit": filter.Limit,
	})
}

// GetTask godoc
// @Summary Получить задачу по ID
// @Tags tasks
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} model.Task
// @Failure 404 {object} model.ErrorResponse
// @Router /api/tasks/{id} [get]
func GetTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Некорректный ID", Code: 400})
		return
	}

	task, err := service.GetTask(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Error: err.Error(), Code: 404})
		return
	}
	c.JSON(http.StatusOK, task)
}

// CreateTask godoc
// @Summary Создать задачу в проекте
// @Tags tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID проекта"
// @Param input body model.TaskInput true "Данные задачи"
// @Success 201 {object} model.Task
// @Failure 400 {object} model.ErrorResponse
// @Router /api/projects/{id}/tasks [post]
func CreateTask(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Некорректный ID", Code: 400})
		return
	}

	var input model.TaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Некорректные данные: " + err.Error(), Code: 400})
		return
	}

	userID := middleware.GetUserID(c)
	task, err := service.CreateTask(uint(projectID), input, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error(), Code: 400})
		return
	}
	c.JSON(http.StatusCreated, task)
}

// UpdateTask godoc
// @Summary Обновить задачу
// @Tags tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Param input body model.TaskInput true "Данные задачи"
// @Success 200 {object} model.Task
// @Failure 400 {object} model.ErrorResponse
// @Router /api/tasks/{id} [put]
func UpdateTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Некорректный ID", Code: 400})
		return
	}

	var input model.TaskUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Некорректные данные: " + err.Error(), Code: 400})
		return
	}

	userID := middleware.GetUserID(c)
	task, err := service.UpdateTask(uint(id), input, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error(), Code: 400})
		return
	}
	c.JSON(http.StatusOK, task)
}

// DeleteTask godoc
// @Summary Удалить задачу
// @Tags tasks
// @Security BearerAuth
// @Param id path int true "ID задачи"
// @Success 200 {object} map[string]string
// @Failure 400 {object} model.ErrorResponse
// @Router /api/tasks/{id} [delete]
func DeleteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Некорректный ID", Code: 400})
		return
	}

	userID := middleware.GetUserID(c)
	if err := service.DeleteTask(uint(id), userID); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error(), Code: 400})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Задача удалена"})
}
