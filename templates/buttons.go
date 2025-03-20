package templates

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

func ButtonTemplateShowMoreOrExit(senderID string) map[string]interface{} {
	return map[string]interface{}{
		"attachment": map[string]interface{}{
			"type": "template",
			"payload": map[string]interface{}{
				"template_type": "button",
				"text":          "Do you want to see more data or exit?",
				"buttons": []map[string]string{
					{"type": "postback", "title": "Show More", "payload": "SHOW_MORE_PAYLOAD"},
					{"type": "postback", "title": "Exit", "payload": "EXIT_PAYLOAD"},
				},
			},
		},
	}
}
