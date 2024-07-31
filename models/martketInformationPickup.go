package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MarketItemPickupInformations struct {
	ID                        uuid.UUID      `json:"id" gorm:"type:varchar(255);primary_key"`
	ItemID                    uuid.UUID      `json:"item_id" gorm:"type:varchar(255);not null"`
	RecipientName             string         `json:"recipient_name" gorm:"type:varchar(255);not null"`
	RecipientPhone            string         `json:"recipient_phone" gorm:"type:varchar(32);not null"`
	Description               string         `json:"description" gorm:"type:varchar(255)"`
	PickupLocationAddress     string         `json:"pickup_location_address" gorm:"type:varchar(255);not null"`
	PickupLocationDescription string         `json:"pickup_location_description" gorm:"type:varchar(255)"`
	ServicePrice              float64        `json:"service_price" gorm:"type:decimal(15,2);not null"`
	DeliveryPrice             float64        `json:"delivery_price" gorm:"type:decimal(15,2);not null"`
	CreatedAt                 time.Time      `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt                 time.Time      `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt                 gorm.DeletedAt `json:"deleted_at" gorm:"type:datetime"`

	MarketItem MarketItems `gorm:"foreignKey:ItemID;references:ID"`
}

func (model *MarketItemPickupInformations) BeforeCreate(tx *gorm.DB) error {
	model.ID = uuid.New()
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	return nil
}

func (model *MarketItemPickupInformations) BeforeUpdate(tx *gorm.DB) error {
	model.UpdatedAt = time.Now()
	return nil
}
