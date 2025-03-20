package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	Object string `json:"object"`
	Entry  []struct {
		ID        string `json:"id"`
		Time      int64  `json:"time"`
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Recipient struct {
				ID string `json:"id"`
			} `json:"recipient"`
			Timestamp int64 `json:"timestamp"`
			Message   struct {
				Mid         string `json:"mid"`
				Text        string `json:"text,omitempty"`
				Attachments []struct {
					Type    string `json:"type"`
					Payload struct {
						URL string `json:"url"`
					} `json:"payload"`
				} `json:"attachments,omitempty"`
			} `json:"message,omitempty"`
			Postback struct {
				Payload string `json:"payload"`
			} `json:"postback,omitempty"`
		} `json:"messaging"`
	} `json:"entry"`
}

type SendMessage struct {
	Recipient struct {
		ID string `json:"id"`
	} `json:"recipient"`
	Message struct {
		Text string `json:"text"`
	} `json:"message"`
}

type Content struct {
	gorm.Model
	StorageContentID uint      `gorm:"index"`
	Timestamp        time.Time `gorm:"not null"`
	Data             string    `gorm:"type:text;not null"`
}

type StorageContent struct {
	gorm.Model
	SenderID    string    `gorm:"size:255;not null"`
	StorageName string    `gorm:"size:255;not null"`
	Contents    []Content `gorm:"foreignKey:StorageContentID"`
	// DeletedAt   gorm.DeletedAt `gorm:"index"`
}
