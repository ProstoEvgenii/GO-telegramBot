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
		if isAdmin(update.Callback.From.Username, whiteList) && update.Callback.Message.Chat.Type == "private" {
			keyboard.HandleSettingsNavigation(update.Callback)
		} else {
			message := models.SendMessage{
				ChatID:    update.Callback.Message.Chat.ID,
				Text:      "Бот выполняет команды только администратора.",
				ParseMode: "None",
			}
			tg.SendMessage(message)
		}

	} else {
		handleMessages(update.Message)
	}
}

func handleMessages(message models.Message) {
	if !isAdmin(message.From.Username, whiteList) && message.Chat.Type != "private" {
		handleUserPublicMessages(message)
	}
	if isAdmin(message.From.Username, whiteList) && message.Chat.Type == "private" {
		handleAdminPrivateMessages(message)
	} else {
		message := models.SendMessage{
			ChatID: message.Chat.ID,
			Text:   "Бот выполняет команды только администратора.",
		}
		tg.SendMessage(message)
		//Добавить отправку сообщения "Вы не являетесь администратором."
	}
}

func handleUserPublicMessages(message models.Message) {
	loweredText := strings.ToLower(message.Text)

	if len(message.Entities) != 0 {
		if !handleEntities(message.Entities, loweredText, whiteList) {
			tg.DeleteMessage(message.Chat.ID, message.MessageID)
		}
	} else if isContainsForbiddenWord(loweredText, forbiddenWords) {
		tg.DeleteMessage(message.Chat.ID, message.MessageID)
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
		userInputHandle(state.InputWord, state.Operation)
		delete(keyboard.UserStates, message.From.Username)
	}
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
func handleBotCommand(message models.Message, messageText string) {
	switch command := isBotCommands(message.Entities, messageText); command {
	case "/settings":
		categoryCommands := []models.InlineKeyboardButton{
			{Text: "Запрещенные слова", CallbackData: "forbidden_words"},
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
func userInputHandle(input, operation string) {
	switch operation {
	case "forbidden_words_add_url":
		log.Printf("Добавил в базу слово %s", input)
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
