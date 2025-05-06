package models

import (
	"time"
)

// Social Media Post represents a task in the system
type SocialMediaLikes struct {
	LikeID    string    `gorm:"type:varchar(75);primaryKey" json:"like_id"`
	PostID    string    `gorm:"type:varchar(75);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"post_id"`
	UserID    string    `gorm:"type:varchar(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user_id"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for Task
func (SocialMediaLikes) TableName() string {
	return "social_media_likes"
}
