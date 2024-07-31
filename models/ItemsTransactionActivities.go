package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MarketItemTransactionActivities struct {
	ID              uuid.UUID      `json:"id" gorm:"type:varchar(255);primary_key"`
	ItemID          uuid.UUID      `json:"item_id" gorm:"type:varchar(255);not null"`
	Status          string         `json:"status" gorm:"type:enum('ON_PROCESS', 'ON_DELIVER', 'FINISHED', 'CANCELLED');not null"`
	TransactionByID uuid.UUID      `json:"transaction_by_id" gorm:"type:varchar(255)"`
	CreatedAt       time.Time      `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"type:datetime"`

	MarketItem MarketItems `gorm:"foreignKey:ItemID;references:ID"`
}

func (model *MarketItemTransactionActivities) BeforeCreate(tx *gorm.DB) error {
	model.ID = uuid.New()
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	return nil
}

func (model *MarketItemTransactionActivities) BeforeUpdate(tx *gorm.DB) error {
	model.UpdatedAt = time.Now()
	return nil
}
