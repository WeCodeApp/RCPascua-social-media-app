package models

import (
	"time"

	"gorm.io/gorm"
)

// Task represents a task in the system
type Task struct {
	ID          string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Title       string         `json:"title" binding:"required" gorm:"type:varchar(255);not null"`
	Description string         `json:"description" gorm:"type:text"`
	Completed   bool           `json:"completed" gorm:"default:false"`
	UserID      string         `json:"user_id" gorm:"type:varchar(36);index;not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName specifies the table name for Task
func (Task) TableName() string {
	return "tasks"
}
