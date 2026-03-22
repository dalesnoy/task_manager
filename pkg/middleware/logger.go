package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware логирует каждый запрос: метод, URL, статус, время ответа
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Обрабатываем запрос
		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		log.Printf("[%s] %s → %d (%v)", method, path, status, duration)
	}
}
