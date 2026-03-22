package repository

import (
	"task-manager/internal/model"
	"task-manager/pkg/database"
)

// GetTagsByTask возвращает теги задачи
func GetTagsByTask(taskID uint) ([]model.Tag, error) {
	var tags []model.Tag
	err := database.DB.Where("task_id = ?", taskID).Find(&tags).Error
	return tags, err
}

// CreateTag создаёт тег
func CreateTag(tag *model.Tag) error {
	return database.DB.Create(tag).Error
}

// DeleteTag удаляет тег
func DeleteTag(id uint) error {
	return database.DB.Delete(&model.Tag{}, id).Error
}

// GetTagByID возвращает тег по ID
func GetTagByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	err := database.DB.First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}
