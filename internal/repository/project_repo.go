package repository

import (
	"task-manager/internal/model"
	"task-manager/pkg/database"
)

// GetProjectsByOwner возвращает все проекты пользователя
func GetProjectsByOwner(ownerID uint) ([]model.Project, error) {
	var projects []model.Project
	err := database.DB.Where("owner_id = ?", ownerID).Find(&projects).Error
	return projects, err
}

// GetProjectByID возвращает проект по ID
func GetProjectByID(id uint) (*model.Project, error) {
	var project model.Project
	err := database.DB.First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// CreateProject создаёт новый проект
func CreateProject(project *model.Project) error {
	return database.DB.Create(project).Error
}

// UpdateProject обновляет проект
func UpdateProject(project *model.Project) error {
	return database.DB.Save(project).Error
}

// DeleteProject удаляет проект и все его задачи
func DeleteProject(id uint) error {
	// Сначала удаляем теги задач этого проекта
	database.DB.Where("task_id IN (SELECT id FROM tasks WHERE project_id = ?)", id).Delete(&model.Tag{})
	// Затем удаляем задачи проекта
	database.DB.Where("project_id = ?", id).Delete(&model.Task{})
	// Удаляем сам проект
	return database.DB.Delete(&model.Project{}, id).Error
}
