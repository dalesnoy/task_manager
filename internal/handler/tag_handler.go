package handler

import (
	"net/http"
	"strconv"

	"task-manager/internal/model"
	"task-manager/internal/service"

	"github.com/gin-gonic/gin"
)

// GetTags godoc
// @Summary Получить теги задачи
// @Tags tags
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {array} model.Tag
// @Router /api/tasks/{id}/tags [get]
func GetTags(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Некорректный ID", Code: 400})
		return
	}

	tags, err := service.GetTags(uint(taskID))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error(), Code: 400})
		return
	}
	c.JSON(http.StatusOK, tags)
}

// CreateTag godoc
// @Summary Добавить тег к задаче
// @Tags tags
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Param input body model.TagInput true "Данные тега"
// @Success 201 {object} model.Tag
// @Failure 400 {object} model.ErrorResponse
// @Router /api/tasks/{id}/tags [post]
func CreateTag(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Некорректный ID", Code: 400})
		return
	}

	var input model.TagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Некорректные данные: " + err.Error(), Code: 400})
		return
	}

	tag, err := service.CreateTag(uint(taskID), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error(), Code: 400})
		return
	}
	c.JSON(http.StatusCreated, tag)
}

// DeleteTag godoc
// @Summary Удалить тег
// @Tags tags
// @Security BearerAuth
// @Param id path int true "ID тега"
// @Success 200 {object} map[string]string
// @Failure 400 {object} model.ErrorResponse
// @Router /api/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Некорректный ID", Code: 400})
		return
	}

	if err := service.DeleteTag(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error(), Code: 400})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Тег удалён"})
}
