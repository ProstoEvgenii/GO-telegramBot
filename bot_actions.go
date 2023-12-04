package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func deleteMessage(chatID, messageID int64) (DeleteMessageResponse, error) {
	route := fmt.Sprintf("deleteMessage?chat_id=%d&message_id=%d", chatID, messageID)

	resp, err := GetToApi(route)
	if err != nil {
		return DeleteMessageResponse{}, fmt.Errorf("error sending request: %s", err)
	}
	var response DeleteMessageResponse

	if err := json.NewDecoder(resp).Decode(&response); err != nil {
		return DeleteMessageResponse{}, fmt.Errorf("error decoding JSON: %s", err)
	}
	log.Println("=b414b5=", response)
	return response, nil
}



func getChatMember(chatID, userID int64) (GetChatMemberResponse, error) {
	route := fmt.Sprintf("getChatMember?chat_id=%d&user_id=%d", chatID, userID)

	resp, err := GetToApi(route)
	if err != nil {
		return GetChatMemberResponse{}, fmt.Errorf("error sending request: %s", err)
	}
	var response GetChatMemberResponse

	if err := json.NewDecoder(resp).Decode(&response); err != nil {
		return GetChatMemberResponse{}, fmt.Errorf("error decoding JSON: %s", err)
	}

	// log.Println("=b414b5=", response)
	return response, nil
}
