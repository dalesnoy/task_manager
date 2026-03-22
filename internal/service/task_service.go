package service

import (
	"errors"

	"task-manager/internal/model"
	"task-manager/internal/repository"
)

var validStatuses = map[string]bool{
	"todo":        true,
	"in_progress": true,
	"done":        true,
}

var validPriorities = map[string]bool{
	"low":    true,
	"medium": true,
	"high":   true,
}

// GetTasks возвращает задачи проекта с фильтрацией
func GetTasks(projectID uint, userID uint, filter model.TaskFilter) ([]model.Task, int64, error) {
	// Проверяем доступ к проекту
	_, err := GetProject(projectID, userID)
	if err != nil {
		return nil, 0, err
	}

	return repository.GetTasksByProject(projectID, filter)
}

// GetTask возвращает задачу по ID
func GetTask(id uint) (*model.Task, error) {
	task, err := repository.GetTaskByID(id)
	if err != nil {
		return nil, errors.New("задача не найдена")
	}
	return task, nil
}

// CreateTask создаёт задачу в проекте
func CreateTask(projectID uint, input model.TaskInput, userID uint) (*model.Task, error) {
	// Проверяем доступ к проекту
	_, err := GetProject(projectID, userID)
	if err != nil {
		return nil, err
	}

	// Устанавливаем значения по умолчанию
	status := input.Status
	if status == "" {
		status = "todo"
	}
	if !validStatuses[status] {
		return nil, errors.New("недопустимый статус: допустимые значения — todo, in_progress, done")
	}

	priority := input.Priority
	if priority == "" {
		priority = "medium"
	}
	if !validPriorities[priority] {
		return nil, errors.New("недопустимый приоритет: допустимые значения — low, medium, high")
	}

	task := model.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      status,
		Priority:    priority,
		Deadline:    input.Deadline,
		ProjectID:   projectID,
		AssigneeID:  input.AssigneeID,
	}

	if err := repository.CreateTask(&task); err != nil {
		return nil, errors.New("ошибка создания задачи")
	}

	return &task, nil
}

// UpdateTask обновляет задачу
func UpdateTask(id uint, input model.TaskUpdateInput, userID uint) (*model.Task, error) {
	task, err := repository.GetTaskByID(id)
	if err != nil {
		return nil, errors.New("задача не найдена")
	}

	// Проверяем доступ к проекту задачи
	_, err = GetProject(task.ProjectID, userID)
	if err != nil {
		return nil, err
	}

	if input.Title != "" {
		task.Title = input.Title
	}
	if input.Description != "" {
		task.Description = input.Description
	}
	if input.Status != "" {
		if !validStatuses[input.Status] {
			return nil, errors.New("недопустимый статус")
		}
		task.Status = input.Status
	}
	if input.Priority != "" {
		if !validPriorities[input.Priority] {
			return nil, errors.New("недопустимый приоритет")
		}
		task.Priority = input.Priority
	}
	if input.Deadline != nil {
		task.Deadline = input.Deadline
	}
	if input.AssigneeID != nil {
		task.AssigneeID = input.AssigneeID
	}

	if err := repository.UpdateTask(task); err != nil {
		return nil, errors.New("ошибка обновления задачи")
	}

	return task, nil
}

// DeleteTask удаляет задачу
func DeleteTask(id uint, userID uint) error {
	task, err := repository.GetTaskByID(id)
	if err != nil {
		return errors.New("задача не найдена")
	}

	_, err = GetProject(task.ProjectID, userID)
	if err != nil {
		return err
	}

	return repository.DeleteTask(id)
}
