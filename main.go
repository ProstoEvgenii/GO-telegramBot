package main

import (
	"GO-chatModeratorTg/models"
	"GO-chatModeratorTg/server"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	// host := "127.0.0.1:80"
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Файл .env не найден")
		// host = ":80"
	}
	Connect()
	filter := bson.M{}

	var forbiddenWords []models.ForbiddenWords
	var whiteList []models.WhiteList
	FindReturnDecoded(filter, "forbiddenWords", &forbiddenWords)
	FindReturnDecoded(filter, "whiteList", &whiteList)
	intervalGetUpdate := 3

	ticker := time.NewTicker(time.Duration(intervalGetUpdate) * time.Second)
	defer ticker.Stop()
	offset := 668578288
	type UserState struct {
		WaitingForInput bool
		InputWord       string
		Operation       string
	}

	var UserStates = make(map[string]UserState) // Хранилище состояний пользователей
	for {
		select {
		case <-ticker.C:
			log.Println("=Выполняюсь каждые=", intervalGetUpdate, "Секунды")

			response, err := getUpdate(offset)
			if err != nil {
				log.Println("=038abf=", err)
			}

			for _, item := range response.Result {
				offset = item.UpdateID + 1
				if item.Callback.ID != "" {

				} else {
					if !isAdmin(item.Message, whiteList) && item.Message.Chat.Type != "private" {
						loweredText := strings.ToLower(item.Message.Text)

						if len(item.Message.Entities) != 0 {
							if !handleEntities(item.Message.Entities, loweredText, whiteList) {
								deleteMessage(item.Message.Chat.ID, item.Message.MessageID)
							}

						} else if isContainsForbiddenWord(loweredText, forbiddenWords) {
							deleteMessage(item.Message.Chat.ID, item.Message.MessageID)
						}
					}
					if isAdmin(item.Message, whiteList) && item.Message.Chat.Type == "private" {
						state, exists := UserStates[item.Message.From.Username]
						//Здесь происходит проверка находится ли пользователь в map Состояния юзеров
						//Если нет - далее,если да проверка  - проверка на ожидание слова
						if !exists {
							loweredText := strings.ToLower(item.Message.Text)
							handleBotCommand(item.Message.Entities, loweredText)
						}
						if state.WaitingForInput {
							state.InputWord = item.Message.Text
							userInputHandle(state.InputWord, state.Operation)
							delete(UserStates, item.Message.From.Username)
						}

					} else {
						//Добавить отправку сообщения "Вы не являетесь администрпатором."
					}
				}
			}
		}
	}
}

func userInputHandle(input, operation string) {
	switch operation {
	case "forbidden_words_add_url":
		log.Printf("Добавил в базу слово %s", input)
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
		keyboard.sendKeybordMessage(2135753546, "Выберите категорию:", [][]models.InlineKeyboardButton{categoryCommands})

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

func getUpdate(offset int) (models.GetUpdates, error) {
	resBody, err := server.GetToApi(fmt.Sprintf("getUpdates?offset=%d", offset))
	if err != nil {
		return models.GetUpdates{}, fmt.Errorf("error fetching data: %s", err)
	}
	defer resBody.Close()

	var response models.GetUpdates
	if err := json.NewDecoder(resBody).Decode(&response); err != nil {
		return models.GetUpdates{}, fmt.Errorf("error decoding JSON: %s", err)
	}

	return response, nil
}

func Start(host string) {
	http.ListenAndServe(host, nil)
}

func sendMenu() error {
	// Формирование меню с описанием
	commands := []models.SetMenu{
		{Command: "/status", Description: "Статус"},
		{Command: "/settings", Description: "Настройки"},
		{Command: "/history", Description: "История"},
	}
	commandsJSON := map[string]interface{}{
		"commands": commands,
	}

	// Формирование сообщения с меню
	requestBody, err := json.Marshal(commandsJSON)
	if err != nil {
		return err
	}
	if res, err := server.PostToApi("setMyCommands", requestBody); err != nil {
		log.Println("=c77107=", res)
		return err
	}
	return nil

}
func writeWords(words []string) {
	for _, forbiddenWord := range words {
		loweredForbiddenWord := strings.ToLower(forbiddenWord)
		filter := bson.M{
			"word": loweredForbiddenWord,
		}
		update := bson.M{"$set": bson.M{
			"word": loweredForbiddenWord,
		}}
		result := InsertIfNotExists(filter, update, "forbiddenWords")
		log.Println("=bcc2f5=", result)
	}

}
func readForbiddenWords(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	return words
}
