package model

import {
	"gorm.io/gorm"
	"github.com/google/uuid"
	"time"
}

type user struct {
	ID UUID `gorm:"type:varchar(255) not null"`
	PhoneNumber string `gorm:"type:varchar(255); unique; not null"`
	Password string `gorm:"type:varchar(32); not null"`
	Role string `gorm:"type:enum('ADMIN','P_SMALL','P_LARGE'.'C_SMALL','C_LARGE'); not null"`
	CreatedAt time.time `gorm:"type:date;default:TIMESTAMP; not null"`
	UpdateAt time.time `gorm:"type:date;default:TIMESTAMP;not null"`
	DeleteAt gorm.DeleteAt `gorm:"type:date"`

}

func CreatedUser (u *user) BeforeCreate(tx * gorm.DB) (err error) {
	u.UUID.new(ID)
	model.user
}