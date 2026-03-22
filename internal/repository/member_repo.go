package repository

import (
	"task-manager/internal/model"
	"task-manager/pkg/database"
)

// AddMember добавляет участника в проект
func AddMember(member *model.ProjectMember) error {
	return database.DB.Create(member).Error
}

// RemoveMember удаляет участника из проекта
func RemoveMember(projectID, userID uint) error {
	return database.DB.Where("project_id = ? AND user_id = ?", projectID, userID).
		Delete(&model.ProjectMember{}).Error
}

// GetProjectMembers возвращает всех участников проекта
func GetProjectMembers(projectID uint) ([]model.ProjectMember, error) {
	var members []model.ProjectMember
	err := database.DB.Preload("User").Where("project_id = ?", projectID).Find(&members).Error
	return members, err
}

// IsProjectMember проверяет, является ли пользователь участником проекта
func IsProjectMember(projectID, userID uint) bool {
	var count int64
	database.DB.Model(&model.ProjectMember{}).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count)
	return count > 0
}
