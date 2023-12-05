package keyboard

import (
	"GO-chatModeratorTg/models"
	"GO-chatModeratorTg/tg"
	"strings"
)

var UserStates = make(map[string]models.UserState) // Хранилище состояний пользователей

func HandleSettings(callbackResult models.CallbackData) {
	path := strings.Split(callbackResult.Data, "_")[0]
	switch path {
	case "home":
		newKeyboard := [][]models.InlineKeyboardButton{
			{
				{Text: "Запрещенные слова", CallbackData: "forbiddenwords"},
				{Text: "WhiteList", CallbackData: "whitelist"},
			},
		}
		tg.EditMessage("Выберите категорию:", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID)
		tg.EditKeybordMessage(callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)

	case "forbiddenwords":
		handleForbiddenWordsNav(callbackResult)
	case "whitelist":
		handleWhitelistNav(callbackResult)
	}
	// log.Println("Received callback:", callbackResult.Data)
}

func handleForbiddenWordsNav(callbackResult models.CallbackData) {
	switch callbackResult.Data {
	case "forbiddenwords":
		newKeyboard := [][]models.InlineKeyboardButton{
			{
				{Text: "Добавить", CallbackData: "forbiddenwords_add"},
				{Text: "Удалить", CallbackData: "forbiddenwords_rm"},
				{Text: "Назад", CallbackData: "home"},
			},
		}
		tg.EditMessage("Выберите действие:", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID)
		tg.EditKeybordMessage(callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)
		delete(UserStates, callbackResult.From.Username)

	case "forbiddenwords_add":
		newKeyboard := [][]models.InlineKeyboardButton{
			{
				{Text: "Назад", CallbackData: "forbiddenwords"},
				{Text: "Отправить еще", CallbackData: "forbiddenwords_add"},
			},
		}
		tg.EditMessage("Отправьте слово, которое нужно добавить.", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID)
		tg.EditKeybordMessage(callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)
		userID := callbackResult.From.Username
		UserStates[userID] = models.UserState{
			WaitingForInput: true,
			Operation:       callbackResult.Data,
			Author:          callbackResult.From.Username,
		}
	case "forbiddenwords_rm":
		newKeyboard := [][]models.InlineKeyboardButton{
			{
				{Text: "Назад", CallbackData: "forbiddenwords"},
				{Text: "Отправить еще", CallbackData: "forbiddenwords_rm"},
			},
		}
		tg.EditMessage("Отправьте слово, которое нужно удалить.", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID)
		tg.EditKeybordMessage(callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)

		userID := callbackResult.From.Username
		UserStates[userID] = models.UserState{
			WaitingForInput: true,
			Operation:       callbackResult.Data,
			Author:          callbackResult.From.Username,
		}
	}
}

func handleWhitelistNav(callbackResult models.CallbackData) {
	switch callbackResult.Data {
	case "whitelist":
		newKeyboard := [][]models.InlineKeyboardButton{
			{
				{Text: "Добавить", CallbackData: "whitelist_add"},
				{Text: "Удалить", CallbackData: "whitelist_rm"},
				{Text: "Назад", CallbackData: "home"},
			},
		}
		tg.EditMessage("Выберите действие:", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID)
		tg.EditKeybordMessage(callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)

	case "whitelist_add":
		newKeyboard := [][]models.InlineKeyboardButton{
			{
				{Text: "Ссылки", CallbackData: "whitelist_add_url"},
				{Text: "Упоминания", CallbackData: "whitelist_add_mention"},
				{Text: "Админ", CallbackData: "whitelist_add_admin"},
				{Text: "Назад", CallbackData: "forbiddenwords"},
			},
		}
		delete(UserStates, callbackResult.From.Username)
		tg.EditMessage("Выберите тип:", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID)
		tg.EditKeybordMessage(callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)

	case "whitelist_add_url":
		newKeyboard := [][]models.InlineKeyboardButton{
			{
				{Text: "Назад", CallbackData: "whitelist_add"},
				{Text: "Отправить еще", CallbackData: "whitelist_add_url"},
			},
		}
		tg.EditMessage("Отправьте ссылку, которую нужно добавить.", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID)
		tg.EditKeybordMessage(callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)

		userID := callbackResult.From.Username
		UserStates[userID] = models.UserState{
			WaitingForInput: true,
			Operation:       callbackResult.Data,
			Type:            "url",
			Author:          callbackResult.From.Username,
		}
	case "whitelist_add_mention":
		newKeyboard := [][]models.InlineKeyboardButton{
			{
				{Text: "Назад", CallbackData: "whitelist_add"},
				{Text: "Отправить еще", CallbackData: "whitelist_add_mention"},
			},
		}
		tg.EditMessage("Отправьте упоминание, которое нужно добавить.", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID)
		tg.EditKeybordMessage(callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)

		userID := callbackResult.From.Username
		UserStates[userID] = models.UserState{
			WaitingForInput: true,
			Operation:       callbackResult.Data,
			Type:            "mention",
			Author:          callbackResult.From.Username,
		}
	case "whitelist_add_admin":
		newKeyboard := [][]models.InlineKeyboardButton{{
			{Text: "Назад", CallbackData: "whitelist_add"},
			{Text: "Отправить еще", CallbackData: "whitelist_add_admin"},
		}}
		tg.EditMessage("Отправьте username администратора, которого нужно добавить.", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID)
		tg.EditKeybordMessage(callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)
		userID := callbackResult.From.Username
		UserStates[userID] = models.UserState{
			WaitingForInput: true,
			Operation:       callbackResult.Data,
			Author:          callbackResult.From.Username,
		}
	case "whitelist_rm":
		newKeyboard := [][]models.InlineKeyboardButton{
			{{Text: "Ссылки", CallbackData: "whitelist_rm_url"},
				{Text: "Упоминания", CallbackData: "whitelist_rm_mention"},
				{Text: "Админ", CallbackData: "whitelist_rm_admin"},
				{Text: "Назад", CallbackData: "forbiddenwords"}},
		}
		delete(UserStates, callbackResult.From.Username)
		tg.EditMessage("Выберите тип:", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID)
		tg.EditKeybordMessage(callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)

	case "whitelist_rm_url":
		newKeyboard := [][]models.InlineKeyboardButton{
			{
				{Text: "Назад", CallbackData: "whitelist_rm"},
				{Text: "Отправить еще", CallbackData: "whitelist_rm_url"},
			},
		}
		tg.EditMessage("Отправьте ссылку, которую нужно удалить.", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID)
		tg.EditKeybordMessage(callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)

		userID := callbackResult.From.Username
		UserStates[userID] = models.UserState{
			WaitingForInput: true,
			Operation:       callbackResult.Data,
			Type:            "url",
			Author:          callbackResult.From.Username,
		}
	case "whitelist_rm_mention":
		newKeyboard := [][]models.InlineKeyboardButton{
			{
				{Text: "Назад", CallbackData: "whitelist_rm"},
				{Text: "Отправить еще", CallbackData: "whitelist_rm_mention"},
			},
		}
		tg.EditMessage("Отправьте упоминание, которое нужно удалить.", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID)
		tg.EditKeybordMessage(callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)

		userID := callbackResult.From.Username
		UserStates[userID] = models.UserState{
			WaitingForInput: true,
			Operation:       callbackResult.Data,
			Type:            "mention",
			Author:          callbackResult.From.Username,
		}
	case "whitelist_rm_admin":
		newKeyboard := [][]models.InlineKeyboardButton{
			{
				{Text: "Назад", CallbackData: "whitelist_rm"},
				{Text: "Отправить еще", CallbackData: "whitelist_rm_admin"},
			},
		}
		tg.EditMessage("Отправьте username администратора, которого нужно удалить.", callbackResult.Message.Chat.ID, callbackResult.Message.MessageID)
		tg.EditKeybordMessage(callbackResult.Message.Chat.ID, callbackResult.Message.MessageID, newKeyboard)
		userID := callbackResult.From.Username
		UserStates[userID] = models.UserState{
			WaitingForInput: true,
			Operation:       callbackResult.Data,
			Type:            "admin",
			Author:          callbackResult.From.Username,
		}
	}
}
