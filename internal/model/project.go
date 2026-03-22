package model

import "time"

type Project struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	OwnerID     uint      `json:"owner_id" gorm:"index;not null"`
	Owner       User      `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	Tasks       []Task    `json:"tasks,omitempty" gorm:"foreignKey:ProjectID"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProjectInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}
