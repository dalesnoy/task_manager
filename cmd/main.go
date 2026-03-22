package main

import (
	"log"
	"os"

	"task-manager/internal/handler"
	"task-manager/internal/model"
	"task-manager/pkg/database"
	"task-manager/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные из .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используем переменные окружения")
	}

	// Подключаемся к базе данных
	database.Connect()

	// Автомиграция — создаём таблицы в БД
	database.DB.AutoMigrate(
		&model.User{},
		&model.Project{},
		&model.Task{},
		&model.Tag{},
	)
	log.Println("Миграция базы данных выполнена")

	// Создаём роутер Gin
	r := gin.New()

	// Подключаем middleware
	r.Use(middleware.LoggerMiddleware()) // логирование запросов
	r.Use(middleware.CORSMiddleware())   // CORS для фронтенда
	r.Use(gin.Recovery())                // восстановление после паники

	// Раздача фронтенда (статические файлы)
	r.Static("/static", "./frontend/build/static")
	r.StaticFile("/", "./frontend/build/index.html")

	// === API маршруты ===
	api := r.Group("/api")

	// Тестовый эндпоинт
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong — сервер работает!"})
	})

	// Auth — регистрация и логин (без JWT)
	auth := api.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
	}

	// Защищённые маршруты (требуют JWT)
	protected := api.Group("")
	protected.Use(middleware.AuthRequired())
	{
		// Проекты
		protected.GET("/projects", handler.GetProjects)
		protected.POST("/projects", handler.CreateProject)
		protected.GET("/projects/:id", handler.GetProject)
		protected.PUT("/projects/:id", handler.UpdateProject)
		protected.DELETE("/projects/:id", handler.DeleteProject)

		// Задачи
		protected.GET("/projects/:id/tasks", handler.GetTasks)
		protected.POST("/projects/:id/tasks", handler.CreateTask)
		protected.GET("/tasks/:id", handler.GetTask)
		protected.PUT("/tasks/:id", handler.UpdateTask)
		protected.DELETE("/tasks/:id", handler.DeleteTask)

		// Теги
		protected.GET("/tasks/:id/tags", handler.GetTags)
		protected.POST("/tasks/:id/tags", handler.CreateTag)
		protected.DELETE("/tags/:id", handler.DeleteTag)
	}

	// Запускаем сервер
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Сервер запущен на порту %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
