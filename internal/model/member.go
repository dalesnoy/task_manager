package model

import "time"

// ProjectMember — участник проекта
type ProjectMember struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProjectID uint      `json:"project_id" gorm:"uniqueIndex:idx_project_user;not null"`
	UserID    uint      `json:"user_id" gorm:"uniqueIndex:idx_project_user;not null"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Role      string    `json:"role" gorm:"default:member;not null"` // owner, member
	CreatedAt time.Time `json:"created_at"`
}

// AddMemberInput — добавление участника по email
type AddMemberInput struct {
	Email string `json:"email" binding:"required,email"`
}
