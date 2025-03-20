package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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

func ListStoragesMessage(senderID string, storages []string) map[string]interface{} {
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
