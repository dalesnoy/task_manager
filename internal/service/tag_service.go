package service

import (
	"errors"

	"task-manager/internal/model"
	"task-manager/internal/repository"
)

// GetTags возвращает теги задачи
func GetTags(taskID uint) ([]model.Tag, error) {
	return repository.GetTagsByTask(taskID)
}

// CreateTag добавляет тег к задаче
func CreateTag(taskID uint, input model.TagInput) (*model.Tag, error) {
	// Проверяем что задача существует
	_, err := repository.GetTaskByID(taskID)
	if err != nil {
		return nil, errors.New("задача не найдена")
	}

	color := input.Color
	if color == "" {
		color = "#cccccc"
	}

	tag := model.Tag{
		Name:   input.Name,
		Color:  color,
		TaskID: taskID,
	}

	if err := repository.CreateTag(&tag); err != nil {
		return nil, errors.New("ошибка создания тега")
	}

	return &tag, nil
}

// DeleteTag удаляет тег
func DeleteTag(id uint) error {
	_, err := repository.GetTagByID(id)
	if err != nil {
		return errors.New("тег не найден")
	}
	return repository.DeleteTag(id)
}
