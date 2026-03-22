package service

import (
	"errors"
	"os"
	"time"

	"task-manager/internal/model"
	"task-manager/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Register регистрирует нового пользователя
func Register(input model.RegisterInput) (*model.User, error) {
	// Проверяем, что email не занят
	existing, _ := repository.GetUserByEmail(input.Email)
	if existing != nil {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	// Хэшируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("ошибка хэширования пароля")
	}

	user := model.User{
		Email:        input.Email,
		PasswordHash: string(hash),
		Name:         input.Name,
	}

	if err := repository.CreateUser(&user); err != nil {
		return nil, errors.New("ошибка создания пользователя")
	}

	return &user, nil
}

// Login проверяет пароль и возвращает JWT токен
func Login(input model.LoginInput) (string, error) {
	user, err := repository.GetUserByEmail(input.Email)
	if err != nil {
		return "", errors.New("неверный email или пароль")
	}

	// Сравниваем хэш пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return "", errors.New("неверный email или пароль")
	}

	// Генерируем JWT токен
	token, err := generateToken(user.ID)
	if err != nil {
		return "", errors.New("ошибка генерации токена")
	}

	return token, nil
}

// generateToken создаёт JWT токен с user_id
func generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(), // токен живёт 72 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}
