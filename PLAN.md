# План разработки таск-менеджера

## Блок 1. Инициализация проекта
- Создание структуры папок (cmd/, internal/, pkg/, frontend/)
- Инициализация Go-модуля (go mod init)
- Создание docker-compose.yml (Go-сервер + PostgreSQL)
- Создание .env с переменными окружения

## Блок 2. Подключение к базе данных
- pkg/database/postgres.go — подключение к PostgreSQL через GORM
- Чтение настроек из .env
- Проверка подключения при старте сервера

## Блок 3. Модели данных (сущности БД)
- internal/model/user.go — User (id, email, password_hash, name, created_at)
- internal/model/project.go — Project (id, title, description, owner_id, created_at)
- internal/model/task.go — Task (id, title, description, status, priority, deadline, project_id, assignee_id)
- internal/model/tag.go — Tag (id, name, color, task_id)
- AutoMigrate при запуске

## Блок 4. Авторизация (Auth)
- internal/repository/user_repo.go — создание и поиск пользователя
- internal/service/auth_service.go — регистрация (bcrypt), логин (JWT)
- internal/handler/auth_handler.go — POST /api/auth/register, POST /api/auth/login
- pkg/middleware/jwt.go — проверка токена

## Блок 5. CRUD для проектов
- internal/repository/project_repo.go
- internal/service/project_service.go
- internal/handler/project_handler.go
- Эндпоинты: GET/POST /api/projects, GET/PUT/DELETE /api/projects/:id

## Блок 6. CRUD для задач
- internal/repository/task_repo.go
- internal/service/task_service.go
- internal/handler/task_handler.go
- Эндпоинты: GET/POST /api/projects/:id/tasks, GET/PUT/DELETE /api/tasks/:id
- Фильтрация: status, priority, page, limit, sort

## Блок 7. Теги для задач
- internal/repository/tag_repo.go
- internal/service/tag_service.go
- internal/handler/tag_handler.go

## Блок 8. Middleware и обработка ошибок
- pkg/middleware/cors.go — CORS
- pkg/middleware/logger.go — логирование запросов
- Единый формат ошибок: {"error": "msg", "code": 404}
- Валидация входящих данных

## Блок 9. Swagger документация
- Аннотации к обработчикам
- Генерация Swagger UI
- Доступ по /swagger/index.html

## Блок 10. Фронтенд
- React-интерфейс
- Страницы: логин, регистрация, список проектов, доска задач
- Подключение к API через fetch/axios

## Блок 11. Docker и деплой
- Dockerfile для Go-сервера
- docker-compose.yml — сервер + PostgreSQL
- Запуск: docker-compose up
