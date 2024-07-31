package model

import{
	"gorn.io/gorm"
	"uuid"
	"time"
}

type user struct {
	ID uudi.ID `gorm:"primarykey" type:"Varchar(255)"; not null`
	PhoneNumber string `gorm:"Phone_number type:password(32)"; not null`
	Password string `gorm:"password" type:"varchar 255"; not null`
	Role string `gorm:"name" type:"enum('ADMIN','P_SMALL','P_LARGE'.'C_SMALL','C_LARGE')"; not null`
	CreatedAt time.time `gorm:"created_at" type:"time.TIMESTAMP"; not null`
	UpdateAt time.time `gorm:"updated_at" type:"time.TIMESTAMP"; not null`
	DeleteAt time.time `gorm: "deleted_at" type:"time.TIMESTAMP";`
}


model.user