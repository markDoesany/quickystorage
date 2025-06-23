package templates

import "fmt"

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
							"title":   "Get Started ",
							"payload": "GET_STARTED_PAYLOAD",
						},
						{
							"type":    "postback",
							"title":   "Help ",
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
							"title":   "Search Storage ",
							"payload": "SEARCH_STORAGE_PAYLOAD",
						},
						{
							"type":    "postback",
							"title":   "Create Storage ",
							"payload": "CREATE_STORAGE_PAYLOAD",
						},
						{
							"type":    "postback",
							"title":   "Remove Storage ",
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

// StorageCarouselTemplate creates a carousel of storage options with up to 3 buttons per card
func StorageCarouselTemplate(senderID string, storages []string, startIndex int) map[string]interface{} {
	const maxItems = 10         // Facebook's limit for carousel items
	const maxButtonsPerCard = 3 // Maximum buttons per card

	// Calculate the end index, making sure not to exceed the slice bounds
	endIndex := startIndex + maxItems
	if endIndex > len(storages) {
		endIndex = len(storages)
	}

	// Create a slice for the current page of storages
	currentStorages := storages[startIndex:endIndex]

	// Calculate how many cards we need (each card will have max 3 buttons)
	cardsNeeded := (len(currentStorages) + maxButtonsPerCard - 1) / maxButtonsPerCard

	// Create elements for the carousel
	elements := make([]map[string]interface{}, 0, cardsNeeded)

	// Process storages in groups of maxButtonsPerCard
	for i := 0; i < len(currentStorages); i += maxButtonsPerCard {
		end := i + maxButtonsPerCard
		if end > len(currentStorages) {
			end = len(currentStorages)
		}

		// Get the current group of storages for this card
		group := currentStorages[i:end]

		// Create buttons for this group
		buttons := make([]map[string]string, 0, len(group))
		for j, storage := range group {
			buttons = append(buttons, map[string]string{
				"type":    "postback",
				"title":   storage,
				"payload": fmt.Sprintf("STORAGE_%d", startIndex+i+j+1), // +1 because storage indices start at 1
			})
		}

		// Create the card
		element := map[string]interface{}{
			"title":    "Select Storage",
			"subtitle": fmt.Sprintf("Page %d/%d", (i/maxButtonsPerCard)+1, cardsNeeded),
			"buttons":  buttons,
		}
		elements = append(elements, element)
	}

	// If there are more items, add a "Next" button to the last element
	if endIndex < len(storages) {
		lastElement := elements[len(elements)-1]
		buttons := lastElement["buttons"].([]map[string]string)
		if len(buttons) < maxButtonsPerCard {
			// Add Next button to existing card if there's space
			buttons = append(buttons, map[string]string{
				"type":    "postback",
				"title":   "Next Page",
				"payload": fmt.Sprintf("STORAGE_PAGE_%d", endIndex),
			})
			lastElement["buttons"] = buttons
		} else {
			// Create a new card for the Next button if no space
			elements = append(elements, map[string]interface{}{
				"title":    "More Options",
				"subtitle": "Continue to next page",
				"buttons": []map[string]string{
					{
						"type":    "postback",
						"title":   "Next Page",
						"payload": fmt.Sprintf("STORAGE_PAGE_%d", endIndex),
					},
				},
			})
		}
	}

	// If this is not the first page, add a "Previous" button to the first element
	if startIndex > 0 {
		firstElement := elements[0]
		buttons := firstElement["buttons"].([]map[string]string)
		if len(buttons) < maxButtonsPerCard {
			// Add Previous button to existing card if there's space
			prevIndex := startIndex - maxItems
			if prevIndex < 0 {
				prevIndex = 0
			}
			buttons = append(buttons, map[string]string{
				"type":    "postback",
				"title":   "Previous Page",
				"payload": fmt.Sprintf("STORAGE_PAGE_%d", prevIndex),
			})
			firstElement["buttons"] = buttons
		} else {
			// Create a new card for the Previous button if no space
			prevIndex := startIndex - maxItems
			if prevIndex < 0 {
				prevIndex = 0
			}
			elements = append([]map[string]interface{}{
				{
					"title":    "Navigation",
					"subtitle": "Return to previous page",
					"buttons": []map[string]string{
						{
							"type":    "postback",
							"title":   "Previous Page",
							"payload": fmt.Sprintf("STORAGE_PAGE_%d", prevIndex),
						},
					},
				},
			}, elements...)
		}
	}

	return map[string]interface{}{
		"recipient": map[string]string{"id": senderID},
		"message": map[string]interface{}{
			"attachment": map[string]interface{}{
				"type": "template",
				"payload": map[string]interface{}{
					"template_type": "generic",
					"elements":      elements,
				},
			},
		},
	}
}
