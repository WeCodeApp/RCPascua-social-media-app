package models

import (
	"time"
)

// Social Media Post represents a task in the system
type SocialMediaPost struct {
	PostID    string    `gorm:"type:varchar(75);primaryKey" json:"post_id"`
	UserID    string    `gorm:"type:varchar(36);not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user_id"`
	PostText  string    `gorm:"type:text;not null;" json:"post_text"`
	PostImage string    `gorm:"type:text;not null;" json:"post_image"`
	Likes     int       `gorm:"type:int;not null" json:"likes"`
	IsLiked   bool      `gorm:"type:bool;not null" json:"is_liked"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for Task
func (SocialMediaPost) TableName() string {
	return "social_media_posts"
}
