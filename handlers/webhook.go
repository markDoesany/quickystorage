package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/markDoesany/quickymessenger/database"
	"github.com/markDoesany/quickymessenger/models"
	"github.com/markDoesany/quickymessenger/services"
	"github.com/markDoesany/quickymessenger/templates"
	"github.com/markDoesany/quickymessenger/utils"
)

var userState = make(map[string]string)
var userStorage = make(map[string][]models.StorageContent)
var mu sync.Mutex

var storage_index int = 0

func InitializeUserStorage(senderID string) {
	var storageContents []models.StorageContent
	err := database.DB.Preload("Contents").Where("sender_id = ?", senderID).Find(&storageContents).Error
	if err != nil {
		log.Fatalf("Failed to load storage contents from database for senderID %s: %v", senderID, err)
	}

	userStorage[senderID] = append(userStorage[senderID], storageContents...)
	log.Printf("User storage initialized from database for senderID %s", senderID)
}

// Webhook handles incoming requests from Messenger
func Webhook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		log.Println("Invalid method: Not GET or POST")
		return
	}

	if r.Method == http.MethodGet {
		verifyToken := r.URL.Query().Get("hub.verify_token")
		if verifyToken != os.Getenv("VERIFY_TOKEN") {
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

	if _, exists := userState[senderID]; !exists {
		InitializeUserStorage(senderID)
		err = services.SendMessage(senderID, templates.ButtonTemplateGetStarted(senderID))
		if err != nil {
			log.Printf("Failed to send message: %v", err)
		}
		userState[senderID] = "waiting_for_get_started"
		return
	}

	if message.Entry[0].Messaging[0].Postback.Payload != "" {
		handlePostbackPayload(senderID, message.Entry[0].Messaging[0].Postback.Payload)
		return
	}

	handleTextInput(senderID, message)
}

func handlePostbackPayload(senderID, payload string) {
	var err error
	switch {
	case payload == "GET_STARTED_PAYLOAD":
		userState[senderID] = "waiting_for_action"
		err = services.SendMessage(senderID, templates.ButtonTemplateMessage(senderID))
	case payload == "SEARCH_STORAGE_PAYLOAD":
		log.Printf("Handling SEARCH_STORAGE_PAYLOAD for senderID: %s", senderID)
		userState[senderID] = "searching"
		storages := getUserStorages(senderID)
		if len(storages) == 0 {
			err = services.SendMessage(senderID, services.TextMessage(senderID, "No storages found."))
			if err != nil {
				log.Fatal(err)
			}
			userState[senderID] = "waiting_for_action"
			err = services.SendMessage(senderID, templates.ButtonTemplateMessage(senderID))
			return
		}
		err = services.SendMessage(senderID, services.ListStoragesMessage(senderID, storages))
		if err != nil {
			log.Fatal(err)
		}
	case strings.HasPrefix(payload, "STORAGE_"):
		// Handle storage selection from carousel or button template
		storageIndexStr := strings.TrimPrefix(payload, "STORAGE_")
		if storageIndexStr == "" {
			err = services.SendMessage(senderID, services.TextMessage(senderID, "Invalid storage selection."))
			break
		}
		index, err := strconv.Atoi(storageIndexStr)
		if err != nil {
			log.Printf("Invalid storage index: %s", storageIndexStr)
			err = services.SendMessage(senderID, services.TextMessage(senderID, "Invalid storage selection."))
			break
		}
		handleStorageSelection(senderID, index)
	case strings.HasPrefix(payload, "STORAGE_PAGE_"):
		// Handle carousel pagination
		pageIndexStr := strings.TrimPrefix(payload, "STORAGE_PAGE_")
		pageIndex, err := strconv.Atoi(pageIndexStr)
		if err != nil {
			log.Printf("Invalid page index: %s", pageIndexStr)
			err = services.SendMessage(senderID, services.TextMessage(senderID, "Invalid page."))
			break
		}
		storages := getUserStorages(senderID)
		err = services.SendMessage(senderID, templates.StorageCarouselTemplate(senderID, storages, pageIndex))
	case payload == "CREATE_STORAGE_PAYLOAD":
		userState[senderID] = "creating"
		err = services.SendMessage(senderID, services.TextMessage(senderID, "Please enter the storage name:"))
	case payload == "REMOVE_STORAGE_PAYLOAD":
		userState[senderID] = "removing"
		storages := getUserStorages(senderID)
		err = services.SendMessage(senderID, services.RemoveListStoragesMessage(senderID, storages))
	case payload == "ADD_DATA_PAYLOAD":
		userState[senderID] = "waiting_for_data"
		err = services.SendMessage(senderID, services.TextMessage(senderID, "Please send a text message (image not yet supported)."))
	case payload == "EXIT_PAYLOAD":
		userState[senderID] = "waiting_for_action"
		err = services.SendMessage(senderID, templates.ButtonTemplateMessage(senderID))
	default:
		if strings.HasPrefix(payload, "REMOVE_STORAGE_") {
			storageIndex := strings.TrimPrefix(payload, "REMOVE_STORAGE_")
			log.Printf("Handling STORAGE Removal for storageIndex: %s", storageIndex)
			index, err := strconv.Atoi(storageIndex)
			if err != nil {
				log.Printf("Invalid storage index: %s", storageIndex)
				return
			}
			storage_index = index - 1
			handleRemoveStorageSelection(senderID, index)
		} else {
			userState[senderID] = "waiting_for_action"
			err = services.SendMessage(senderID, services.TextMessage(senderID, "Invalid selection. Please choose an option."))
			if err == nil {
				err = services.SendMessage(senderID, templates.ButtonTemplateMessage(senderID))
			}
		}
	}
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func handleTextInput(senderID string, message models.Message) {
	var err error
	state, exists := userState[senderID]
	if exists {
		switch state {
		case "creating":
			storageName := message.Entry[0].Messaging[0].Message.Text
			log.Printf("Creating storage with name: %s for senderID: %s", storageName, senderID)
			userStorage[senderID] = append(userStorage[senderID], models.StorageContent{StorageName: storageName, Contents: []models.Content{}})
			userState[senderID] = "storing_data"
			storages := getUserStorages(senderID)
			storage_index = len(storages) - 1
			err = services.SendMessage(senderID, services.TextMessage(senderID, "Storage created: *"+storageName+"*"))
			if err == nil {
				err = services.SendMessage(senderID, templates.ButtonTemplateAddOrExit(senderID))
			}
		case "storing_data":
			userState[senderID] = "waiting_for_data"
			err = services.SendMessage(senderID, services.TextMessage(senderID, "Please send a text message or an image."))
		case "waiting_for_data":
			var data string
			// if len(message.Entry[0].Messaging[0].Message.Attachments) > 0 {
			// 	attachment := message.Entry[0].Messaging[0].Message.Attachments[0]
			// 	if attachment.Type == "image" {
			// 		data = attachment.Payload.URL
			// 		log.Printf("Received image URL: %s", data)
			// 	} else {
			// 		log.Printf("Unsupported attachment type: %s", attachment.Type)
			// 		err = services.SendMessage(senderID, services.TextMessage(senderID, "Unsupported attachment type. Please send an image."))
			// 		break
			// 	}
			// } else {
			// }
			data = message.Entry[0].Messaging[0].Message.Text
			timestamp := time.Now()
			log.Printf("Storing data: %s with timestamp: %s for senderID: %s", data, timestamp, senderID)
			err = database.StoreDataInDB(senderID, userStorage[senderID][storage_index].StorageName, timestamp, data)
			if err != nil {
				log.Printf("Failed to store data in database: %v", err)
				break
			}
			userState[senderID] = "storing_data"
			err = services.SendMessage(senderID, services.TextMessage(senderID, "Data stored: "+data+"."))
			if err == nil {
				err = services.SendMessage(senderID, templates.ButtonTemplateAddOrExit(senderID))
			}
		case "searching":
			storageName := message.Entry[0].Messaging[0].Message.Text
			log.Printf("Searching storage with name: %s for senderID: %s", storageName, senderID)
		default:
			userState[senderID] = "waiting_for_action"
			err = services.SendMessage(senderID, services.TextMessage(senderID, "I didn't understand that. Click a button to proceed."))
			if err == nil {
				err = services.SendMessage(senderID, templates.ButtonTemplateMessage(senderID))
			}
		}
	} else {
		err = services.SendMessage(senderID, services.TextMessage(senderID, "Click 'Get Started' to begin."))
	}
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func handleStorageSelection(senderID string, index int) error {
	storages, exists := userStorage[senderID]
	if !exists || len(storages) == 0 {
		log.Printf("No storages found for senderID: %s", senderID)
		return services.SendMessage(senderID, services.TextMessage(senderID, "You don't have any storages."))
	}

	storage := storages[index-1]
	log.Printf("Retrieving storage: %s for senderID: %s", storage.StorageName, senderID)

	contents, err := database.GetStorageData(senderID, storage.StorageName)
	if err != nil {
		log.Printf("Failed to get storage content: %v", err)
		return err
	}
	if len(contents) == 0 {
		services.SendMessage(senderID, services.TextMessage(senderID, "No data found in storage: _"+storage.StorageName+"_"))
	}

	log.Printf("Storage contents for %s: %v", storage.StorageName, contents)
	err = services.SendMessage(senderID, services.TextMessage(senderID, "Storage Name: "+storage.StorageName))
	if err != nil {
		log.Printf("Failed to send storage content: %v", err)
		return err
	}
	for _, content := range contents {
		responseMessage := "Timestamp:\n" + utils.FormatTimestamp(content.Timestamp) + "\n\nData:\n" + content.Data
		err := services.SendMessage(senderID, services.TextMessage(senderID, responseMessage))
		if err != nil {
			log.Printf("Failed to send storage content: %v", err)
			return err
		}
	}

	userState[senderID] = "waiting_for_action"
	err = services.SendMessage(senderID, templates.ButtonTemplateAddOrExit(senderID))

	return err
}

func handleRemoveStorageSelection(senderID string, index int) error {
	storages, exists := userStorage[senderID]
	if !exists || len(storages) == 0 {
		log.Printf("No storages found for senderID: %s", senderID)
		return services.SendMessage(senderID, services.TextMessage(senderID, "You don't have any storages."))
	}

	storage_ := storages[index-1]
	log.Printf("Retrieving storage: %s for senderID: %s", storage_.StorageName, senderID)

	for i, storage := range userStorage[senderID] {
		if storage.StorageName == storage_.StorageName {
			userStorage[senderID] = append(userStorage[senderID][:i], userStorage[senderID][i+1:]...)
			break
		}
	}
	userState[senderID] = "waiting_for_action"
	err := services.SendMessage(senderID, templates.ButtonTemplateMessage(senderID))

	return err
}

func getUserStorages(senderID string) []string {
	storages := []string{}
	for _, storage := range userStorage[senderID] {
		storages = append(storages, storage.StorageName)
	}
	return storages
}
