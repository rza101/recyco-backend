package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ForumPost struct {
	ID           uuid.UUID      `json:"id" gorm:"type:varchar(255);primary_key"`
	Title        string         `json:"title" gorm:"type:varchar(255);not null"`
	Description  string         `json:"description" gorm:"type:text;not null"`
	ThumbnailUrl string         `json:"thumbnail_url" gorm:"type:varchar(2048);null"`
	CreatedBy    uuid.UUID      `json:"created_by" gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time      `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"type:datetime"`

	CreatedByUser User             `json:"created_by_user" gorm:"foreignKey:CreatedBy;references:ID"`
	PostId        []ForumPostReply `json:"forum_post_replies" gorm:"foreignKey:post_id;references:ID"`
}

func (model *ForumPost) BeforeCreate(tx *gorm.DB) error {
	model.ID = uuid.New()
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	return nil
}

func (model *ForumPost) BeforeUpdate(tx *gorm.DB) error {
	model.UpdatedAt = time.Now()
	return nil
}
