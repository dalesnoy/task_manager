package middleware

import (
	"net/http"
	"os"
	"strings"

	"task-manager/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthRequired проверяет JWT токен в заголовке Authorization
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем заголовок Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "Требуется авторизация",
				Code:  http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Формат: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "Неверный формат токена",
				Code:  http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		secret := os.Getenv("JWT_SECRET")

		// Парсим и проверяем токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "Недействительный токен",
				Code:  http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Извлекаем user_id из токена и сохраняем в контекст
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "Ошибка чтения токена",
				Code:  http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		userID := uint(claims["user_id"].(float64))
		c.Set("user_id", userID)

		// Сохраняем флаг админа
		if isAdmin, ok := claims["is_admin"].(bool); ok {
			c.Set("is_admin", isAdmin)
		} else {
			c.Set("is_admin", false)
		}

		c.Next()
	}
}

// GetUserID извлекает user_id из контекста Gin
func GetUserID(c *gin.Context) uint {
	userID, _ := c.Get("user_id")
	return userID.(uint)
}

// AdminRequired проверяет, что пользователь — администратор
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			c.JSON(http.StatusForbidden, model.ErrorResponse{
				Error: "Доступ только для администраторов",
				Code:  http.StatusForbidden,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
