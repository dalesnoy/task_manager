package service

import (
	"errors"

	"task-manager/internal/model"
	"task-manager/internal/repository"
)

// AddMember добавляет участника в проект по email
func AddMember(projectID uint, input model.AddMemberInput, ownerID uint) (*model.ProjectMember, error) {
	// Только владелец может добавлять участников
	project, err := repository.GetProjectByID(projectID)
	if err != nil {
		return nil, errors.New("проект не найден")
	}
	if project.OwnerID != ownerID {
		return nil, errors.New("только владелец может добавлять участников")
	}

	// Ищем пользователя по email
	user, err := repository.GetUserByEmail(input.Email)
	if err != nil {
		return nil, errors.New("пользователь с таким email не найден")
	}

	// Нельзя добавить самого себя
	if user.ID == ownerID {
		return nil, errors.New("вы уже являетесь владельцем проекта")
	}

	// Проверяем, не добавлен ли уже
	if repository.IsProjectMember(projectID, user.ID) {
		return nil, errors.New("пользователь уже является участником проекта")
	}

	member := model.ProjectMember{
		ProjectID: projectID,
		UserID:    user.ID,
		Role:      "member",
	}

	if err := repository.AddMember(&member); err != nil {
		return nil, errors.New("ошибка добавления участника")
	}

	// Подгружаем данные пользователя для ответа
	member.User = *user

	return &member, nil
}

// RemoveMember удаляет участника из проекта
func RemoveMember(projectID uint, memberUserID uint, ownerID uint) error {
	// Только владелец может удалять участников
	project, err := repository.GetProjectByID(projectID)
	if err != nil {
		return errors.New("проект не найден")
	}
	if project.OwnerID != ownerID {
		return errors.New("только владелец может удалять участников")
	}

	if !repository.IsProjectMember(projectID, memberUserID) {
		return errors.New("пользователь не является участником проекта")
	}

	return repository.RemoveMember(projectID, memberUserID)
}

// GetMembers возвращает всех участников проекта (включая владельца первым)
func GetMembers(projectID uint, userID uint) ([]model.ProjectMember, error) {
	// Доступ есть у владельца и участников
	if !HasProjectAccess(projectID, userID) {
		return nil, errors.New("нет доступа к этому проекту")
	}

	// Получаем проект, чтобы добавить владельца
	project, err := repository.GetProjectByID(projectID)
	if err != nil {
		return nil, errors.New("проект не найден")
	}

	// Загружаем данные владельца
	owner, err := repository.GetUserByID(project.OwnerID)
	if err != nil {
		return nil, errors.New("ошибка загрузки владельца")
	}

	// Владелец первым в списке
	result := []model.ProjectMember{
		{
			ProjectID: projectID,
			UserID:    owner.ID,
			User:      *owner,
			Role:      "owner",
		},
	}

	// Остальные участники
	members, err := repository.GetProjectMembers(projectID)
	if err != nil {
		return nil, err
	}

	result = append(result, members...)
	return result, nil
}

// HasProjectAccess проверяет, есть ли у пользователя доступ к проекту (владелец или участник)
func HasProjectAccess(projectID uint, userID uint) bool {
	project, err := repository.GetProjectByID(projectID)
	if err != nil {
		return false
	}
	if project.OwnerID == userID {
		return true
	}
	return repository.IsProjectMember(projectID, userID)
}
