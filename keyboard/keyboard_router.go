package keyboard

import (
	"GO-chatModeratorTg/models"
	"GO-chatModeratorTg/server"
	"encoding/json"
	"log"
)

func HandleForbidenWordNavigation(callbackResult models.CallbackData) {
	switch callback := callbackResult.Data; callback {
	case "home":
		newKeyboard := [][]models.InlineKeyboardButton{
			{
				{Text: "Запрещенные слова", CallbackData: "forbidden_words"},
				{Text: "WhiteList", CallbackData: "whitelist"},
			},
		}
		editKeybordMessage("Выберите категорию:", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)

	case "forbidden_words":
		newKeyboard := [][]models.InlineKeyboardButton{
			{
				{Text: "Добавить", CallbackData: "forbidden_words_add"},
				{Text: "Удалить", CallbackData: "forbidden_words_remove"},
				{Text: "Назад", CallbackData: "home"},
			},
		}
		editKeybordMessage("Выберите действие:", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)
	case "forbidden_words_add":
		newKeyboard := [][]models.InlineKeyboardButton{
			{
				{Text: "Ссылки", CallbackData: "forbidden_words_add_url"},
				{Text: "Упоминания", CallbackData: "forbidden_words_add_mention"},
				{Text: "Админ", CallbackData: "forbidden_words_add_admin"},
				{Text: "Назад", CallbackData: "forbidden_words"},
			},
		}
		delete(UserStates, callbackResult.From.Username)
		editKeybordMessage("Выберите тип:", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)
	case "forbidden_words_add_url":
		newKeyboard := [][]models.InlineKeyboardButton{
			{
				{Text: "Назад", CallbackData: "forbidden_words_add"},
			},
		}

		editKeybordMessage("Выберите тип:", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)
		userID := callbackResult.From.Username
		word := ""
		main.UserStates[userID] = models.UserState{
			WaitingForInput: true,
			InputWord:       word,
			Operation:       callback,
		}
	}

	log.Println("Received callback:", callbackResult.Data)

}

func editKeybordMessage(text string, chatID int64, messageID int64, keyboard [][]models.InlineKeyboardButton) {
	message := models.SendMessage{
		ChatID:    chatID,
		Text:      text,
		MessageID: messageID,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		},
	}
	messageJSON, _ := json.Marshal(message)
	if _, err := server.PostToApi("editMessageReplyMarkup", messageJSON); err != nil {
		log.Println("=52a1d9=", err)
	}
}
func sendKeybordMessage(chatID int64, text string, keyboard [][]models.InlineKeyboardButton) {
	message := models.SendMessage{
		ChatID: chatID,
		Text:   text,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		},
	}

	messageJSON, _ := json.Marshal(message)
	if _, err := server.PostToApi("sendMessage", messageJSON); err != nil {
		log.Println("=52a1d9=", err)
	}
}
