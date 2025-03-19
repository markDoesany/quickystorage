package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

const (
	GRAPHQL_URL  = "https://graph.facebook.com/v2.6"
	ACCESS_TOKEN = "EAAIqTJOUImEBO0tvOusaJzZBLCh0ASuSNZC5LUhp66SavfelN5WOxbZAc6343iZCwRU9MATWoFHQBiZCdS1MReZByWJ6uKlep6PJ5x2BNVZBnXqP06gHWZCtJR3j4FrxiNYop3ZAjDeL2VLtH8EuZAjdaaxeMzxvUlw1tZCGEpOt4FudgzdxCpCJY1SosOa1cjhkgEeRAZDZD"
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

	url := fmt.Sprintf("%s/%s?access_token=%s", GRAPHQL_URL, "me/messages", ACCESS_TOKEN)
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
			"payload": "STORAGE_" + string(i+1),
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

func ButtonTemplateGetStarted(senderID string) map[string]interface{} {
	return map[string]interface{}{
		"recipient": map[string]string{"id": senderID},
		"message": map[string]interface{}{
			"attachment": map[string]interface{}{
				"type": "template",
				"payload": map[string]interface{}{
					"template_type": "button",
					"text":          "What would you like to do?",
					"buttons": []map[string]string{
						{
							"type":    "postback",
							"title":   "Get Started 🚀",
							"payload": "GET_STARTED_PAYLOAD",
						},
						{
							"type":    "postback",
							"title":   "Help ❓",
							"payload": "HELP_PAYLOAD",
						},
					},
				},
			},
		},
	}
}

func ButtonTemplateMessage(senderID string) map[string]interface{} {
	return map[string]interface{}{
		"recipient": map[string]string{"id": senderID},
		"message": map[string]interface{}{
			"attachment": map[string]interface{}{
				"type": "template",
				"payload": map[string]interface{}{
					"template_type": "button",
					"text":          "What would you like to do?",
					"buttons": []map[string]string{
						{
							"type":    "postback",
							"title":   "Search Storage 🔍",
							"payload": "SEARCH_STORAGE_PAYLOAD",
						},
						{
							"type":    "postback",
							"title":   "Create Storage ✍️",
							"payload": "CREATE_STORAGE_PAYLOAD",
						},
						{
							"type":    "postback",
							"title":   "Remove Storage ❌",
							"payload": "REMOVE_STORAGE_PAYLOAD",
						},
					},
				},
			},
		},
	}
}

func ButtonTemplateAddOrExit(senderID string) map[string]interface{} {
	return map[string]interface{}{
		"recipient": map[string]string{"id": senderID},
		"message": map[string]interface{}{
			"attachment": map[string]interface{}{
				"type": "template",
				"payload": map[string]interface{}{
					"template_type": "button",
					"text":          "Do you want to add more data or exit?",
					"buttons": []map[string]string{
						{
							"type":    "postback",
							"title":   "Add Data",
							"payload": "ADD_DATA_PAYLOAD",
						},
						{
							"type":    "postback",
							"title":   "Exit",
							"payload": "EXIT_PAYLOAD",
						},
					},
				},
			},
		},
	}
}
