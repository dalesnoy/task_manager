package model

type Tag struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name" gorm:"not null"`
	Color  string `json:"color" gorm:"default:#cccccc"`
	TaskID uint   `json:"task_id" gorm:"index;not null"`
}

type TagInput struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color"`
}
