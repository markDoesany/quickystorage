package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Storage struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}

type Message struct {
	gorm.Model
	StorageID uint   `gorm:"index"`
	Text      string `gorm:"not null"`
}

var DB *gorm.DB

func InitDB() {
	dsn := "root:#GoodKredit.com@tcp(127.0.0.1:3306)/quickystorage_bot?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	DB.AutoMigrate(&Storage{}, &Message{})
	fmt.Println("Database connected and migrated successfully!")
}
