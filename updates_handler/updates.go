package updates

import (
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
	db.ForbiddenWords(filter, "whiteList", &whiteList)
}

func UpdatesHandler(update models.Result) {
	if update.Callback.ID != "" {
		keyboard.HandleForbidenWordNavigation(update.Callback)
	} else {
		handleMessages(update.Message)
	}
}
func handleMessages(message models.Message) {
	if !isAdmin(message, whiteList) && message.Chat.Type != "private" {
		loweredText := strings.ToLower(message.Text)

		if len(message.Entities) != 0 {
			if !handleEntities(message.Entities, loweredText, whiteList) {
				tg.DeleteMessage(message.Chat.ID, message.MessageID)
			}

		} else if isContainsForbiddenWord(loweredText, forbiddenWords) {
			tg.DeleteMessage(message.Chat.ID, message.MessageID)
		}
	}
	if isAdmin(message, whiteList) && message.Chat.Type == "private" {
		state, exists := keyboard.UserStates[message.From.Username]
		//Здесь происходит проверка находится ли пользователь в map Состояния юзеров
		//Если нет - далее,если да проверка  - проверка на ожидание слова
		if !exists {
			loweredText := strings.ToLower(message.Text)
			handleBotCommand(message.Entities, loweredText)
		}
		if state.WaitingForInput {
			state.InputWord = message.Text
			userInputHandle(state.InputWord, state.Operation)
			delete(keyboard.UserStates, message.From.Username)
		}

	} else {
		//Добавить отправку сообщения "Вы не являетесь администрпатором."
	}
}
func isAdmin(message models.Message, whiteList []models.WhiteList) bool {
	user := message.From.Username
	adminsMap := make(map[string]bool)
	for _, item := range whiteList {
		if item.Type == "admin" {
			adminsMap[item.Content] = true
		}
	}
	return adminsMap[user]
}
func handleBotCommand(entities []models.Entities, messageText string) {
	switch command := isBotCommands(entities, messageText); command {
	case "/settings":
		categoryCommands := []models.InlineKeyboardButton{
			{Text: "Запрещенные слова", CallbackData: "forbidden_words"},
			{Text: "WhiteList", CallbackData: "whitelist"},
		}
		keyboard.SendKeybordMessage(2135753546, "Выберите категорию:", [][]models.InlineKeyboardButton{categoryCommands})

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
