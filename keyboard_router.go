package main

import "log"

func HandleForbidenWordNavigation(callbackResult CallbackData) {
	switch callback := callbackResult.Data; callback {
	case "home":
		newKeyboard := [][]InlineKeyboardButton{
			{
				{Text: "Запрещенные слова", CallbackData: "forbidden_words"},
				{Text: "WhiteList", CallbackData: "whitelist"},
			},
		}
		editKeybordMessage("Выберите категорию:", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)

	case "forbidden_words":
		newKeyboard := [][]InlineKeyboardButton{
			{
				{Text: "Добавить", CallbackData: "forbidden_words_add"},
				{Text: "Удалить", CallbackData: "forbidden_words_remove"},
				{Text: "Назад", CallbackData: "home"},
			},
		}
		editKeybordMessage("Выберите действие:", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)
	case "forbidden_words_add":
		newKeyboard := [][]InlineKeyboardButton{
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
		newKeyboard := [][]InlineKeyboardButton{
			{
				{Text: "Назад", CallbackData: "forbidden_words_add"},
			},
		}

		editKeybordMessage("Выберите тип:", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)
		userID := callbackResult.From.Username
		word := ""
		UserStates[userID] = UserState{
			WaitingForInput: true,
			InputWord:       word,
			Operation:       callback,
		}
	}

	log.Println("Received callback:", callbackResult.Data)

}
