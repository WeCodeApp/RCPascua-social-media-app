package models

import (
	"time"
)

// Social Media Post represents a task in the system
type SocialMediaComments struct {
	CommentID   string    `gorm:"type:varchar(75);primaryKey" json:"comment_id"`
	PostID      string    `gorm:"type:varchar(75);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"post_id"`
	UserID      string    `gorm:"type:varchar(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user_id"`
	CommentText string    `gorm:"type:text;not null;" json:"comment_text"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for Task
func (SocialMediaComments) TableName() string {
	return "social_media_comments"
}
