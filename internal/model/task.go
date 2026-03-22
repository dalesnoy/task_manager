package model

import "time"

type Task struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	Status      string     `json:"status" gorm:"default:todo;not null"`       // todo, in_progress, done
	Priority    string     `json:"priority" gorm:"default:medium;not null"`   // low, medium, high
	Deadline    *time.Time `json:"deadline"`
	ProjectID   uint       `json:"project_id" gorm:"index;not null"`
	AssigneeID  *uint      `json:"assignee_id" gorm:"index"`
	Assignee    *User      `json:"assignee,omitempty" gorm:"foreignKey:AssigneeID"`
	IsPrivate   bool       `json:"is_private" gorm:"default:false"`       // приватная задача — видна только владельцу
	CreatorID   uint       `json:"creator_id" gorm:"index"`
	Tags        []Tag      `json:"tags,omitempty" gorm:"foreignKey:TaskID"`
	CreatedAt   time.Time  `json:"created_at"`
}

// TaskInput — для создания задачи (title обязателен)
type TaskInput struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	Deadline    *time.Time `json:"deadline"`
	AssigneeID  *uint      `json:"assignee_id"`
	IsPrivate   bool       `json:"is_private"`
}

// TaskUpdateInput — для обновления задачи (все поля опциональны)
type TaskUpdateInput struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	Deadline    *time.Time `json:"deadline"`
	AssigneeID  *uint      `json:"assignee_id"`
	IsPrivate   *bool      `json:"is_private"`
}

// TaskFilter — параметры фильтрации и пагинации
type TaskFilter struct {
	Status   string `form:"status"`
	Priority string `form:"priority"`
	Sort     string `form:"sort"`
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
}
