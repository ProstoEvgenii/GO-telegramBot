package main

import (
	"GO-chatModeratorTg/models"
	"GO-chatModeratorTg/server"
	"encoding/json"
	"fmt"
	"log"
)

func deleteMessage(chatID, messageID int64) (models.DeleteMessageResponse, error) {
	route := fmt.Sprintf("deleteMessage?chat_id=%d&message_id=%d", chatID, messageID)

	resp, err := server.GetToApi(route)
	if err != nil {
		return models.DeleteMessageResponse{}, fmt.Errorf("error sending request: %s", err)
	}
	var response models.DeleteMessageResponse

	if err := json.NewDecoder(resp).Decode(&response); err != nil {
		return models.DeleteMessageResponse{}, fmt.Errorf("error decoding JSON: %s", err)
	}
	log.Println("=b414b5=", response)
	return response, nil
}

func getChatMember(chatID, userID int64) (models.GetChatMemberResponse, error) {
	route := fmt.Sprintf("getChatMember?chat_id=%d&user_id=%d", chatID, userID)

	resp, err := server.GetToApi(route)
	if err != nil {
		return models.GetChatMemberResponse{}, fmt.Errorf("error sending request: %s", err)
	}
	var response models.GetChatMemberResponse

	if err := json.NewDecoder(resp).Decode(&response); err != nil {
		return models.GetChatMemberResponse{}, fmt.Errorf("error decoding JSON: %s", err)
	}

	// log.Println("=b414b5=", response)
	return response, nil
}
