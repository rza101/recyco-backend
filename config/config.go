package config

import (
	"recyco/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/recyco?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	database.AutoMigrate(
		&models.User{},
		&models.Article{},
		&models.Community{},
		&models.MarketItems{},
		&models.MarketItemTransactionActivities{},
		&models.MarketItemPickupInformations{},
		&models.ForumPost{},
		&models.ForumPostReply{},
		&models.TreatmentLocation{},
	)
	DB = database
}
