package updates

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

var whiteList []models.WhiteList
var forbiddenWords []models.ForbiddenWords

func GetWhiteListAndForbiddeWords() {
	filter := bson.M{}
	db.FindReturnDecoded(filter, "forbiddenWords", &forbiddenWords)
	db.FindReturnDecoded(filter, "whiteList", &whiteList)
}

func UpdatesHandler(update models.Result) {
	if update.Callback.ID != "" {
		handleCallback(update.Callback)

	} else {
		handleMessages(update.Message)
	}
}

func handleMessages(message models.Message) {
	if !isAdmin(message.From.Username, whiteList) && message.Chat.Type != "private" {
		//Сообщение пользователя в чат.
		handleUserPublicMessages(message)
	}
	if isAdmin(message.From.Username, whiteList) && message.Chat.Type == "private" {
		//Сообщения админа в личку.
		handleAdminPrivateMessages(message)
	}
	// tg.SendMessage(message)
	//Добавить отправку сообщения "Вы не являетесь администратором."
	// }
}

func handleUserPublicMessages(message models.Message) {
	loweredText := strings.ToLower(message.Text)

	if len(message.Entities) != 0 {
		if !handleEntities(message.Entities, loweredText, whiteList) {
			//Проверка сущностей используемых в сообщении(ссылки,упоминания)
			tg.DeleteMessage(message.Chat.ID, message.MessageID)
		}
	} else if isContainsForbiddenWord(loweredText, forbiddenWords) {
		//Проверка на запрещенные слова
		tg.DeleteMessage(message.Chat.ID, message.MessageID)
	}
}

func handleCallback(callBack models.CallbackData) {
	if isAdmin(callBack.From.Username, whiteList) && callBack.Message.Chat.Type == "private" {
		keyboard.HandleSettingsCallback(callBack)
	} else if !isAdmin(callBack.From.Username, whiteList) && callBack.Message.Chat.Type == "private" {
		message := models.SendMessage{
			ChatID:    callBack.Message.Chat.ID,
			Text:      "Бот выполняет команды только администратора.",
			ParseMode: "None",
		}
		tg.SendMessage(message)
	}
}
func handleAdminPrivateMessages(message models.Message) {
	state, exists := keyboard.UserStates[message.From.Username]
	//Здесь происходит проверка находится ли пользователь в map Состояния юзеров
	//Если нет - далее,если да проверка  - проверка на ожидание слова
	if !exists {
		loweredText := strings.ToLower(message.Text)
		handleBotCommand(message, loweredText)
	}
	if state.WaitingForInput {
		state.InputWord = message.Text
		userInputHandle(state, message.Chat.ID)
		delete(keyboard.UserStates, message.From.Username)
	}
}

func handleBotCommand(message models.Message, messageText string) {
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

	case "/status":
		fmt.Println(command, "is great!")
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
func userInputHandle(userInput models.UserState, ChatID int64) {
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
			message.Text = fmt.Sprintf(`Добавил в базу слово "%s".`, userInput.InputWord)
			tg.SendMessage(message)
		} else {
			message.Text = fmt.Sprintf(`Добавил в базу слово "%s".`, userInput.InputWord)
			tg.SendMessage(message)
		}
	case "forbiddenwords_rm":
		filter := bson.M{
			"word": userInput.InputWord,
		}
		if db.DeleteDocument(filter, "forbiddenWords") {
			// log.Printf(`Удалил слово %s из базы.`, userInput.InputWord)
			message.Text = fmt.Sprintf(`Удалил слово "%s" из базы.`, userInput.InputWord)
			tg.SendMessage(message)
		} else {
			// log.Printf(`Cлово "%s" отсутствует в бд.`, userInput.InputWord)
			message.Text = fmt.Sprintf(`Cлово "%s" отсутствует в бд.`, userInput.InputWord)
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
		log.Printf(`Удалено  %s %s из базы .`, userInput.Type, userInput.InputWord)
		message.Text = fmt.Sprintf(`Удалено %s из whiteList.`, userInput.InputWord)
		tg.SendMessage(message)
	} else {
		message.Text = fmt.Sprintf(`%s "%s" отсутствует в whiteList.`, userInput.Type, userInput.InputWord)
		tg.SendMessage(message)
	}
}

func handleEntities(entities []models.Entities, messageText string, whiteList []models.WhiteList) bool {

	for _, entity := range entities {
		if entity.Type == "mention" {
			//Обрабатывать разрешенные упоминания
			start := entity.Offset
			stop := start + entity.Length
			mention := messageText[start:stop]
			for _, item := range whiteList {
				if item.Type == "mention" && item.Content == mention {
					return true
				}
			}
			return false
		}
		if entity.Type == "url" {
			start := entity.Offset
			stop := start + entity.Length
			url := messageText[start:stop]
			//Обрабатывать разрешенные ссылки
			for _, item := range whiteList {
				if item.Type == entity.Type && strings.Contains(url, item.Content) {
					return true
				}
			}
			return false

		}
		if entity.Type == "text_link" {
			//Обрабатывать разрешенные ссылки
			for _, item := range whiteList {
				if item.Type == "url" && strings.Contains(entity.URL, item.Content) {
					return true
				}
			}
			return false
		}

	}
	return false
}
func isContainsForbiddenWord(message string, forbiddenWords []models.ForbiddenWords) bool {
	message = regexp.MustCompile("[^a-zA-Zа-яА-Я0-9\\s]+").ReplaceAllString(message, "")
	messageWords := strings.Fields(message)

	forbiddenMap := make(map[string]bool)
	for _, fword := range forbiddenWords {
		forbiddenMap[strings.ToLower(fword.Word)] = true
	}

	for _, messageWord := range messageWords {
		if forbiddenMap[messageWord] {
			log.Printf("=Сообщение будет удалено из-за слова %s", messageWord)
			return true

		}
	}
	return false
}
func isAdmin(user string, whiteList []models.WhiteList) bool {
	adminsMap := make(map[string]bool)

	for _, item := range whiteList {
		if item.Type == "admin" {
			adminsMap[item.Content] = true
		}
	}
	return adminsMap[user]
}
