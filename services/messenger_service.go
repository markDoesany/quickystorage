package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/markDoesany/quickymessenger/templates"
)

type SendRequestFunc func(senderID string) (interface{}, error)

func SendMessage(senderID string, message map[string]interface{}) error {
	if message == nil {
		return errors.New("message can't be empty")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling request: %w", err)
	}

	url := fmt.Sprintf("%s/%s?access_token=%s", os.Getenv("GRAPHQL_URL"), "me/messages", os.Getenv("ACCESS_TOKEN"))
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	log.Printf("message sent successfully?\n%#v", res)

	return nil
}

func TextMessage(senderID, text string) map[string]interface{} {
	return map[string]interface{}{
		"recipient": map[string]string{"id": senderID},
		"message": map[string]string{
			"text": text,
		},
	}
}

func PayloadMessage(senderID, payload string) SendRequestFunc {
	return func(senderID string) (interface{}, error) {
		if payload == "" {
			return nil, errors.New("payload can't be empty")
		}
		return map[string]interface{}{
			"recipient": map[string]string{"id": senderID},
			"postback":  map[string]string{"payload": payload},
		}, nil
	}
}

// ListStoragesMessage creates a message with a list of storages
// Uses carousel template when there are more than 3 storages
func ListStoragesMessage(senderID string, storages []string) map[string]interface{} {
	// If 3 or fewer storages, use button template for simplicity
	if len(storages) <= 3 {
		buttons := make([]map[string]string, len(storages))
		for i, storage := range storages {
			buttons[i] = map[string]string{
				"type":    "postback",
				"title":   storage,
				"payload": "STORAGE_" + strconv.Itoa(i+1),
			}
		}

		return map[string]interface{}{
			"recipient": map[string]string{"id": senderID},
			"message": map[string]interface{}{
				"attachment": map[string]interface{}{
					"type": "template",
					"payload": map[string]interface{}{
						"template_type": "button",
						"text":          "Select a storage:",
						"buttons":       buttons,
					},
				},
			},
		}
	}

	// For more than 3 storages, use carousel template
	return templates.StorageCarouselTemplate(senderID, storages, 0)
}

func RemoveListStoragesMessage(senderID string, storages []string) map[string]interface{} {
	buttons := make([]map[string]string, len(storages))
	for i, storage := range storages {
		buttons[i] = map[string]string{
			"type":    "postback",
			"title":   storage,
			"payload": "REMOVE_STORAGE_" + strconv.Itoa(i+1),
		}
	}

	return map[string]interface{}{
		"recipient": map[string]string{"id": senderID},
		"message": map[string]interface{}{
			"attachment": map[string]interface{}{
				"type": "template",
				"payload": map[string]interface{}{
					"template_type": "button",
					"text":          "Select a storage:",
					"buttons":       buttons,
				},
			},
		},
	}
}

func RemoveListStorages(senderID string, storages []string) map[string]interface{} {
	buttons := make([]map[string]string, len(storages))
	for i, storage := range storages {
		buttons[i] = map[string]string{
			"type":    "postback",
			"title":   storage,
			"payload": "REMOVESTORAGE_" + strconv.Itoa(i+1),
		}
	}

	return map[string]interface{}{
		"recipient": map[string]string{"id": senderID},
		"message": map[string]interface{}{
			"attachment": map[string]interface{}{
				"type": "template",
				"payload": map[string]interface{}{
					"template_type": "button",
					"text":          "Select a storage:",
					"buttons":       buttons,
				},
			},
		},
	}
}

// SetupPersistentMenu configures the persistent menu for the Messenger bot
func SetupPersistentMenu() error {
	log.Println("Setting up persistent menu...")

	payload := map[string]interface{}{
		"persistent_menu": []map[string]interface{}{
			{
				"locale":                  "default",
				"composer_input_disabled": false,
				"call_to_actions": []map[string]interface{}{
					// My Account
					{"type": "postback", "title": "View Profile", "payload": "VIEW_PROFILE"},
					{"type": "postback", "title": "Create Storage", "payload": "CREATE_STORAGE"},
					{"type": "postback", "title": "Billing Statement", "payload": "BILLING_STATEMENT"},
					{"type": "postback", "title": "Payment History", "payload": "PAYMENT_HISTORY"},
					{"type": "postback", "title": "Update Info", "payload": "UPDATE_INFO"},
					// Requests & Concerns
					{"type": "postback", "title": "Maintenance Request", "payload": "MAINTENANCE_REQUEST"},
					{"type": "postback", "title": "Book Amenity", "payload": "BOOK_AMENITY"},
					{"type": "postback", "title": "Visitor Pass", "payload": "VISITOR_PASS"},
					{"type": "postback", "title": "Complaint", "payload": "COMPLAINT"},
					{"type": "postback", "title": "Feedback", "payload": "FEEDBACK"},
					// Community
					{"type": "postback", "title": "Announcements", "payload": "ANNOUNCEMENTS"},
					{"type": "postback", "title": "Events", "payload": "EVENTS"},
					{"type": "postback", "title": "Buy/Sell Board", "payload": "BUY_SELL_BOARD"},
					{"type": "postback", "title": "Contact Admin", "payload": "CONTACT_ADMIN"},
					{"type": "postback", "title": "Community Guidelines", "payload": "COMMUNITY_GUIDELINES"},
				},
			},
		},
	}

	// Log the payload structure for debugging
	payloadJSON, _ := json.MarshalIndent(payload, "", "  ")
	log.Printf("Payload to be sent: %s\n", payloadJSON)

	// Get the access token from environment
	accessToken := os.Getenv("ACCESS_TOKEN")
	if accessToken == "" {
		return fmt.Errorf("ACCESS_TOKEN environment variable is not set")
	}

	// Construct the URL
	url := fmt.Sprintf("https://graph.facebook.com/v19.0/me/messenger_profile?access_token=%s", accessToken)
	log.Printf("Sending request to: %s\n", url)

	// Marshal the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v\n", err)
		return fmt.Errorf("error marshaling payload: %v", err)
	}

	// Create and send the request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read and log the response body
	body, _ := io.ReadAll(resp.Body)
	log.Printf("Response status: %s\n", resp.Status)
	log.Printf("Response body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error setting up messenger profile (status %d): %s", resp.StatusCode, string(body))
	}

	log.Println("Messenger profile set up successfully!")
	return nil
}
