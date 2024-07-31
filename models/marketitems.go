package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MarketItems struct {
	ID            uuid.UUID      `json:"id" gorm:"type:varchar(255);primary_key"`
	Name          string         `json:"name" gorm:"type:varchar(255);not null"`
	Price         float64        `json:"price" gorm:"type:decimal(15,2);not null"`
	Weight        float64        `json:"weight" gorm:"type:decimal(10,2);not null"`
	ItemScale     string         `json:"item_scale" gorm:"type:enum('SMALL', 'LARGE');not null"`
	Description   string         `json:"description" gorm:"type:text"`
	ThumbnailUrl  string         `json:"thumbnail_url" gorm:"type:varchar(2048)"`
	PostedBy      uuid.UUID      `json:"posted_by" gorm:"type:varchar(255);not null"`
	OrderedBy     *uuid.UUID     `json:"ordered_by" gorm:"type:varchar(255)"`
	CreatedAt     time.Time      `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"type:datetime"`
	PostedByUser  User           `json:"posted_by_user" gorm:"foreignKey:PostedBy;references:ID"`
	OrderedByUser User           `json:"ordered_by_user" gorm:"foreignKey:OrderedBy;references:ID"`

	TransactionActivities []MarketItemTransactionActivities `json:"transaction_activities" gorm:"foreignKey:ItemID;references:ID"`
	PickupInformations    []MarketItemPickupInformations    `json:"pickup_informations" gorm:"foreignKey:ItemID;references:ID"`
}

func (model *MarketItems) BeforeCreate(tx *gorm.DB) error {
	model.ID = uuid.New()
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	return nil
}

func (model *MarketItems) BeforeUpdate(tx *gorm.DB) error {
	model.UpdatedAt = time.Now()
	return nil
}
