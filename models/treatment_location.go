package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TreatmentLocation struct {
	ID          uuid.UUID      `json:"id" gorm:"type:varchar(255);primary_key"`
	Title       string         `json:"title" gorm:"type:varchar(255);not null"`
	Address     string         `json:"address" gorm:"type:varchar(255);not null"`
	Lat         float64        `json:"lat" gorm:"type:float;not null"`
	Lon         float64        `json:"lon" gorm:"type:float;not null"`
	Description string         `json:"description" gorm:"type:varchar(255)"`
	CreatedAt   time.Time      `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"type:datetime"`
}

func (model *TreatmentLocation) BeforeCreate(tx *gorm.DB) error {
	model.ID = uuid.New()
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	return nil
}

func (model *TreatmentLocation) BeforeUpdate(tx *gorm.DB) error {
	model.UpdatedAt = time.Now()
	return nil
}
