package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Community struct {
	ID           uuid.UUID      `json:"id" gorm:"type:varchar(255);primary_key"`
	Name         string         `json:"name" gorm:"type:varchar(255);not null"`
	Description  string         `json:"description" gorm:"type:text; not null"`
	Community    string         `json:"community_url" gorm:"type:varchar(255);not null"`
	ThumbnailUrl string         `json:"thumbnail_url" gorm:"type:varchar(2048);no null"`
	CreatedAt    time.Time      `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"type:datetime"`
}

func (model *Community) BeforeCreate(tx *gorm.DB) error {
	model.ID = uuid.New()
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	return nil
}

func (model *Community) BeforeUpdate(tx *gorm.DB) error {
	model.UpdatedAt = time.Now()
	return nil
}
