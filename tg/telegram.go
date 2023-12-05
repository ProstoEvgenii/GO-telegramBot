package tg

import (
	"GO-chatModeratorTg/models"
	"GO-chatModeratorTg/server"
	"encoding/json"
	"fmt"
	"log"
)

func DeleteMessage(chatID, messageID int64) (models.DeleteMessageResponse, error) {
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

func SendMessage(message models.SendMessage) {
	messageJSON, _ := json.Marshal(message)
	if _, err := server.PostToApi("sendMessage", messageJSON); err != nil {
		log.Println("=52a1d9=", err)
	}
}
