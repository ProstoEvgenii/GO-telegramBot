package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func deleteMessage(chatID, messageID int) (Delete, error) {
	route := fmt.Sprintf("deleteMessage?chat_id=%d&message_id=%d", chatID, messageID)

	resp, err := GetToApi(route)
	if err != nil {
		return Delete{}, fmt.Errorf("error sending request: %s", err)
	}
	var response Delete

	if err := json.NewDecoder(resp).Decode(&response); err != nil {
		return Delete{}, fmt.Errorf("error decoding JSON: %s", err)
	}
	if !response.Result {
		return Delete{}, fmt.Errorf("error deleting message: %t", response.Result)
	}
	log.Println("=b414b5=", response)
	return response, nil
}
