package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/markDoesany/quickymessenger/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DSN")
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&models.StorageContent{}, &models.Content{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	fmt.Println("Database connected and migrated successfully!")
}

func StoreDataInDB(senderID, storageName string, timestamp time.Time, data string) error {
	var storageContent models.StorageContent
	err := DB.Where("sender_id = ? AND storage_name = ?", senderID, storageName).First(&storageContent).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			storageContent = models.StorageContent{
				SenderID:    senderID,
				StorageName: storageName,
			}
			if err := DB.Create(&storageContent).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	content := models.Content{
		StorageContentID: storageContent.ID,
		Timestamp:        timestamp,
		Data:             data,
	}
	return DB.Create(&content).Error
}

func GetStorageData(senderID, storageName string) ([]models.Content, error) {
	var storageContent models.StorageContent
	err := DB.Where("sender_id = ? AND storage_name = ?", senderID, storageName).First(&storageContent).Error
	if err != nil {
		return nil, err
	}

	var contents []models.Content
	err = DB.Where("storage_content_id = ?", storageContent.ID).Find(&contents).Error
	if err != nil {
		return nil, err
	}
	log.Print(contents)
	return contents, nil
}
