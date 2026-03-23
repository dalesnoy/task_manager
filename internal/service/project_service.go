package service

import (
	"errors"

	"task-manager/internal/model"
	"task-manager/internal/repository"
)

// GetProjects возвращает все проекты пользователя (свои + где участник)
func GetProjects(userID uint) ([]model.Project, error) {
	return repository.GetProjectsByUser(userID)
}

// GetProject возвращает проект по ID с проверкой прав
func GetProject(id uint, userID uint) (*model.Project, error) {
	project, err := repository.GetProjectByID(id)
	if err != nil {
		return nil, errors.New("проект не найден")
	}

	if project.OwnerID != userID && !repository.IsProjectMember(id, userID) {
		return nil, errors.New("нет доступа к этому проекту")
	}

	return project, nil
}

// CreateProject создаёт новый проект
func CreateProject(input model.ProjectInput, userID uint) (*model.Project, error) {
	project := model.Project{
		Title:       input.Title,
		Description: input.Description,
		OwnerID:     userID,
	}

	if err := repository.CreateProject(&project); err != nil {
		return nil, errors.New("ошибка создания проекта")
	}

	return &project, nil
}

// UpdateProject обновляет проект
func UpdateProject(id uint, input model.ProjectInput, userID uint) (*model.Project, error) {
	project, err := GetProject(id, userID)
	if err != nil {
		return nil, err
	}

	project.Title = input.Title
	project.Description = input.Description

	if err := repository.UpdateProject(project); err != nil {
		return nil, errors.New("ошибка обновления проекта")
	}

	return project, nil
}

// DeleteProject удаляет проект
func DeleteProject(id uint, userID uint) error {
	_, err := GetProject(id, userID)
	if err != nil {
		return err
	}

	return repository.DeleteProject(id)
}
