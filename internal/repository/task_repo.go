package repository

import (
	"task-manager/internal/model"
	"task-manager/pkg/database"
)

// GetTasksByProject возвращает задачи проекта с фильтрацией и пагинацией
// userID используется для фильтрации приватных задач
func GetTasksByProject(projectID uint, userID uint, filter model.TaskFilter) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64

	// Показываем публичные задачи + приватные задачи текущего пользователя
	query := database.DB.Where("project_id = ? AND (is_private = false OR creator_id = ?)", projectID, userID)

	// Фильтрация по статусу
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// Фильтрация по приоритету
	if filter.Priority != "" {
		query = query.Where("priority = ?", filter.Priority)
	}

	// Считаем общее количество (для пагинации)
	query.Model(&model.Task{}).Count(&total)

	// Сортировка
	switch filter.Sort {
	case "deadline":
		query = query.Order("deadline ASC")
	case "priority":
		query = query.Order("priority DESC")
	case "created_at":
		query = query.Order("created_at DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// Пагинация
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit

	err := query.Preload("Tags").Limit(filter.Limit).Offset(offset).Find(&tasks).Error
	return tasks, total, err
}

// GetTaskByID возвращает задачу по ID
func GetTaskByID(id uint) (*model.Task, error) {
	var task model.Task
	err := database.DB.Preload("Tags").First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// CreateTask создаёт задачу
func CreateTask(task *model.Task) error {
	return database.DB.Create(task).Error
}

// UpdateTask обновляет задачу
func UpdateTask(task *model.Task) error {
	return database.DB.Save(task).Error
}

// DeleteTask удаляет задачу и её теги
func DeleteTask(id uint) error {
	database.DB.Where("task_id = ?", id).Delete(&model.Tag{})
	return database.DB.Delete(&model.Task{}, id).Error
}
