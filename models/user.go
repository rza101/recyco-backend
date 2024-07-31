package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                 uuid.UUID        `json:"id" gorm:"type:varchar(255);primary_key"`
	PhoneNumber        string           `json:"phone_number" gorm:"type:varchar(255);not null;unique"`
	Password           string           `json:"password" gorm:"type:varchar(255);not null"`
	Name               string           `json:"name" gorm:"type:varchar(255);not null"`
	Role               string           `json:"role" gorm:"type:enum('ADMIN', 'P_SMALL', 'P_LARGE', 'C_SMALL', 'C_LARGE')"`
	CreatedAt          time.Time        `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time        `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt          gorm.DeletedAt   `json:"deleted_at" gorm:"type:datetime"`
	Articles           []Article        `json:"articles" gorm:"foreignKey:CreatedBy;references:ID"`
	MarketItemsPosted  []MarketItems    `json:"market_items_posted" gorm:"foreignKey:PostedBy;references:ID"`
	MarketItemsOrdered []MarketItems    `json:"market_items_ordered" gorm:"foreignKey:OrderedBy;references:ID"`
	ForumPosts         []ForumPost      `json:"forum_posts" gorm:"foreignKey:CreatedBy;references:ID"`
	ForumPostReplies   []ForumPostReply `json:"forum_post_replies" gorm:"foreignKey:RepliedBy;references:ID"`
}

func (model *User) BeforeCreate(tx *gorm.DB) error {
	model.ID = uuid.New()
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	return nil
}

func (model *User) BeforeUpdate(tx *gorm.DB) error {
	model.UpdatedAt = time.Now()
	return nil
}
