package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/markDoesany/quickymessenger/models"
	"github.com/markDoesany/quickymessenger/services"
)

const VERIFY_TOKEN = "12345"

// Temporary in-memory storage to track user interactions
var userState = make(map[string]string) // senderID -> state
var mu sync.Mutex                       // Mutex to handle concurrency

// Webhook handles incoming requests from Messenger
func Webhook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		log.Println("Invalid method: Not GET or POST")
		return
	}

	if r.Method == http.MethodGet {
		verifyToken := r.URL.Query().Get("hub.verify_token")
		if verifyToken != VERIFY_TOKEN {
			log.Println("Invalid verification token")
			return
		}
		if _, err := w.Write([]byte(r.URL.Query().Get("hub.challenge"))); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read body: %v", err)
		return
	}

	var message models.Message
	if err := json.Unmarshal(body, &message); err != nil {
		log.Printf("Failed to unmarshal body: %v", err)
		return
	}

	if len(message.Entry) == 0 || len(message.Entry[0].Messaging) == 0 {
		log.Println("Invalid message format")
		return
	}

	senderID := message.Entry[0].Messaging[0].Sender.ID

	mu.Lock()
	defer mu.Unlock()

	// Send the "Get Started" button template initially
	if _, exists := userState[senderID]; !exists {
		err = services.SendMessage(senderID, services.ButtonTemplateGetStarted(senderID))
		if err != nil {
			log.Printf("Failed to send message: %v", err)
		}
		userState[senderID] = "waiting_for_get_started"
		return
	}

	// Handle postback payloads
	if message.Entry[0].Messaging[0].Postback.Payload != "" {
		payload := message.Entry[0].Messaging[0].Postback.Payload

		switch payload {
		case "GET_STARTED_PAYLOAD":
			userState[senderID] = "waiting_for_action"
			err = services.SendMessage(senderID, services.ButtonTemplateMessage(senderID))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			return

		case "SEARCH_STORAGE_PAYLOAD":
			userState[senderID] = "searching"
			// Send list of storages (dummy data for now)
			storages := []string{"Storage 1", "Storage 2", "Storage 3"}
			err = services.SendMessage(senderID, services.ListStoragesMessage(senderID, storages))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			return

		case "CREATE_STORAGE_PAYLOAD":
			userState[senderID] = "creating"
			err = services.SendMessage(senderID, services.TextMessage(senderID, "Please enter the storage name:"))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			return

		case "REMOVE_STORAGE_PAYLOAD":
			userState[senderID] = "removing"
			// Send list of storages (dummy data for now)
			storages := []string{"Storage 1", "Storage 2", "Storage 3"}
			err = services.SendMessage(senderID, services.ListStoragesMessage(senderID, storages))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			return

		case "ADD_DATA_PAYLOAD":
			userState[senderID] = "waiting_for_data"
			err = services.SendMessage(senderID, services.TextMessage(senderID, "Please send a text message or an image."))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			return

		case "EXIT_PAYLOAD":
			userState[senderID] = "waiting_for_action"
			err = services.SendMessage(senderID, services.ButtonTemplateMessage(senderID))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			return

		default:
			userState[senderID] = "waiting_for_action"
			err = services.SendMessage(senderID, services.TextMessage(senderID, "Invalid selection. Please choose an option."))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			err = services.SendMessage(senderID, services.ButtonTemplateMessage(senderID))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}
		}
	}

	// Process text-based user input (for follow-ups after button clicks)
	state, exists := userState[senderID]
	if exists {
		switch state {
		case "creating":
			// Store the storage name and ask for data
			storageName := message.Entry[0].Messaging[0].Message.Text
			userState[senderID] = "storing_data"
			err = services.SendMessage(senderID, services.TextMessage(senderID, "Storage created: "+storageName+". Do you want to add data or exit?"))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			err = services.SendMessage(senderID, services.ButtonTemplateAddOrExit(senderID))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}

		case "storing_data":
			// Ask for text or image data
			userState[senderID] = "waiting_for_data"
			err = services.SendMessage(senderID, services.TextMessage(senderID, "Please send a text message or an image."))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}

		case "waiting_for_data":
			// Store the data and ask for more or exit
			data := message.Entry[0].Messaging[0].Message.Text // Handle image data separately
			userState[senderID] = "storing_data"
			err = services.SendMessage(senderID, services.TextMessage(senderID, "Data stored: "+data+". Do you want to add more data or exit?"))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			err = services.SendMessage(senderID, services.ButtonTemplateAddOrExit(senderID))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}

		case "searching":
			// Handle storage search and return data
			storageIndex := message.Entry[0].Messaging[0].Message.Text
			// Retrieve and send storage data (dummy data for now)
			storageData := "Data for storage " + storageIndex
			err = services.SendMessage(senderID, services.TextMessage(senderID, storageData))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			userState[senderID] = "waiting_for_action"
			err = services.SendMessage(senderID, services.ButtonTemplateMessage(senderID))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}

		case "removing":
			// Handle storage removal
			storageIndex := message.Entry[0].Messaging[0].Message.Text
			// Remove storage (dummy action for now)
			err = services.SendMessage(senderID, services.TextMessage(senderID, "Storage "+storageIndex+" removed."))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			userState[senderID] = "waiting_for_action"
			err = services.SendMessage(senderID, services.ButtonTemplateMessage(senderID))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}

		default:
			userState[senderID] = "waiting_for_action"
			err = services.SendMessage(senderID, services.TextMessage(senderID, "I didn't understand that. Click a button to proceed."))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			err = services.SendMessage(senderID, services.ButtonTemplateMessage(senderID))
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}
		}
	} else {
		err = services.SendMessage(senderID, services.TextMessage(senderID, "Click 'Get Started' to begin."))
		if err != nil {
			log.Printf("Failed to send message: %v", err)
		}
	}
}
