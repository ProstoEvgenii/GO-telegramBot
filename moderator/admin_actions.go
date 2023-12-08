package moderator

import (
	"GO-chatModeratorTg/db"
	"GO-chatModeratorTg/keyboard"
	"GO-chatModeratorTg/models"
	"GO-chatModeratorTg/tg"
	"fmt"
	"log"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func handleAdminPrivateMessages(message models.Message) {
	state, exists := keyboard.UserStates[message.From.Username]
	//Здесь происходит проверка находится ли админ в карте,которая попалняется,
	//когда админ хочет записать что-то в базу.
	//Если нет - далее,если да проверка  - проверка на ожидание слова
	if !exists {
		loweredText := strings.ToLower(message.Text)
		handleAdminBotCommand(message, loweredText)
	}
	if state.WaitingForInput {
		state.InputWord = message.Text
		adminInputHandle(state, message.Chat.ID)
		delete(keyboard.UserStates, message.From.Username)
	}
}

func handleAdminBotCommand(message models.Message, messageText string) {
	switch command := isBotCommands(message.Entities, messageText); command {
	case "/settings":
		categoryCommands := []models.InlineKeyboardButton{
			{Text: "Запрещенные слова", CallbackData: "forbiddenwords"},
			{Text: "WhiteList", CallbackData: "whitelist"},
		}
		messageKeybord := models.SendMessage{
			ChatID: message.Chat.ID,
			Text:   "Выберите категорию:",
			ReplyMarkup: &models.InlineKeyboardMarkup{
				InlineKeyboard: [][]models.InlineKeyboardButton{categoryCommands},
			},
		}
		tg.SendMessage(messageKeybord)

	case "/site":
		newMessage := models.SendMessage{
			ChatID:    message.Chat.ID,
			ParseMode: "None",
			Text:      "https://tt1.cryptodev.store/\nlogin: dima\npassword: !SuperPassDlyaDimy!",
		}
		tg.SendMessage(newMessage)
	case "/start":
		filter := bson.M{
			"content": message.From.Username,
			"type":    "admin",
		}
		log.Println("=46c309=", message.Chat.ID)
		update := bson.M{"$set": bson.M{
			"chatID": message.Chat.ID,
		}}
		db.InsertIfNotExists(filter, update, "whiteList", false)
		GetWhiteListAndForbiddeWords()

	default:

	}
}

func isBotCommands(entities []models.Entities, messageText string) string {
	for _, entity := range entities {
		if entity.Type == "bot_command" {
			//Обрабатывать разрешенные упоминания
			start := entity.Offset
			stop := start + entity.Length
			command := messageText[start:stop]
			return command
		}
	}
	return ""
}

func adminInputHandle(userInput models.UserState, ChatID int64) {
	message := models.SendMessage{
		ChatID:    ChatID,
		ParseMode: "None",
	}
	switch userInput.Operation {
	case "forbiddenwords_add":
		filter := bson.M{
			"word": userInput.InputWord,
		}
		update := bson.M{"$set": bson.M{
			"word":    userInput.InputWord,
			"addedBy": userInput.Author,
		}}
		if db.InsertIfNotExists(filter, update, "forbiddenWords", true) {
			message.Text = fmt.Sprintf(`Добавил в запрещенное слово "%s".`, userInput.InputWord)
			tg.SendMessage(message)
			GetWhiteListAndForbiddeWords()
		} else {
			message.Text = fmt.Sprintf(`Слово "%s" уже запрещено.`, userInput.InputWord)
			tg.SendMessage(message)
		}
	case "forbiddenwords_rm":
		filter := bson.M{
			"word": userInput.InputWord,
		}
		if db.DeleteDocument(filter, "forbiddenWords") {
			message.Text = fmt.Sprintf(`Слово "%s" больше не запрещено.`, userInput.InputWord)
			tg.SendMessage(message)
			GetWhiteListAndForbiddeWords()
		} else {
			message.Text = fmt.Sprintf(`Cлово "%s" отсутствует в базе запрещенных слов.`, userInput.InputWord)
			tg.SendMessage(message)
		}
	case "whitelist_add_url":
		addToWhiteList(userInput, ChatID)
	case "whitelist_add_mention":
		addToWhiteList(userInput, ChatID)
	case "whitelist_add_admin":
		addToWhiteList(userInput, ChatID)
	case "whitelist_rm_url":
		removeFromWhiteList(userInput, ChatID)
	case "whitelist_rm_mention":
		removeFromWhiteList(userInput, ChatID)
	case "whitelist_rm_admin":
		removeFromWhiteList(userInput, ChatID)
	}
}
func addToWhiteList(userInput models.UserState, ChatID int64) {
	if userInput.Type != "url" {
		re := regexp.MustCompile(`@`)
		userInput.InputWord = re.ReplaceAllString(userInput.InputWord, "")
	}
	// log.Println("=e2ee53=", userInput)
	filter := bson.M{
		"content": userInput.InputWord,
		"type":    userInput.Type,
	}
	update := bson.M{"$set": bson.M{
		"content": userInput.InputWord,
		"type":    userInput.Type,
		"addedBy": userInput.Author,
	}}
	message := models.SendMessage{
		ChatID:    ChatID,
		ParseMode: "None",
	}
	if db.InsertIfNotExists(filter, update, "whiteList", true) {
		//Добавить вызов функции для записи в лог кем и что сделано
		message.Text = fmt.Sprintf(`Добавлен %s "%s" в WhiteList.`, userInput.Type, userInput.InputWord)
		tg.SendMessage(message)
		GetWhiteListAndForbiddeWords()
	} else {
		message.Text = fmt.Sprintf(`%s "%s" уже в WhiteList.`, userInput.Type, userInput.InputWord)
		tg.SendMessage(message)
	}
}

func removeFromWhiteList(userInput models.UserState, ChatID int64) {
	if userInput.Type != "url" {
		re := regexp.MustCompile(`@`)
		userInput.InputWord = re.ReplaceAllString(userInput.InputWord, "")
	}
	filter := bson.M{
		"content": userInput.InputWord,
		"type":    userInput.Type,
	}
	message := models.SendMessage{
		ChatID:    ChatID,
		ParseMode: "None",
	}
	if db.DeleteDocument(filter, "whiteList") {
		// log.Printf(`Удалено  %s %s из базы .`, userInput.Type, userInput.InputWord)
		message.Text = fmt.Sprintf(`Удалено %s из whiteList.`, userInput.InputWord)
		tg.SendMessage(message)
		GetWhiteListAndForbiddeWords()
	} else {
		message.Text = fmt.Sprintf(`%s "%s" отсутствует в whiteList.`, userInput.Type, userInput.InputWord)
		tg.SendMessage(message)
	}
}
