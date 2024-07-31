package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ForumPostReply struct {
	ID          uuid.UUID      `json:"id" gorm:"type:varchar(255);primary_key"`
	PostID      uuid.UUID      `json:"post_id" gorm:"type:varchar(255);not null"`
	Description string         `json:"description" gorm:"type:text;not null"`
	RepliedBy   uuid.UUID      `json:"replied_by" gorm:"type:varchar(255);not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"type:datetime"`

	RepliedByUser User `json:"replied_by_user" gorm:"foreignKey:RepliedBy;references:ID"`
}

func (model *ForumPostReply) BeforeCreate(tx *gorm.DB) error {
	model.ID = uuid.New()
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	return nil
}

func (model *ForumPostReply) BeforeUpdate(tx *gorm.DB) error {
	model.UpdatedAt = time.Now()
	return nil
}
