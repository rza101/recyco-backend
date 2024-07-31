package config

import {
	"recyo/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
}

func config () {
	dsn := "user:root;pass@tcp(127.0.0.1:3306)/dbname:recyco?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

db.AutoMigrate(
	&user,
	&article,
)