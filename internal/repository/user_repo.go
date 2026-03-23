package repository

import (
	"task-manager/internal/model"
	"task-manager/pkg/database"
)

// CreateUser создаёт нового пользователя в БД
func CreateUser(user *model.User) error {
	return database.DB.Create(user).Error
}

// GetUserByEmail находит пользователя по email
func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID находит пользователя по ID
func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAllUsers возвращает всех пользователей
func GetAllUsers() ([]model.User, error) {
	var users []model.User
	err := database.DB.Order("created_at DESC").Find(&users).Error
	return users, err
}
