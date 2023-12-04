package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

type SetMenu struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type SendMessage struct {
	ChatID      int64                 `json:"chat_id"`
	MessageID   int64                 `json:"message_id"`
	Text        string                `json:"text"`
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}
type UserState struct {
	WaitingInput bool // Флаг ожидания ввода
}

func main() {
	// host := "127.0.0.1:80"
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Файл .env не найден")
		// host = ":80"
	}
	Connect()
	filter := bson.M{}

	var forbiddenWords []ForbiddenWords
	var whiteList []WhiteList
	FindReturnDecoded(filter, "forbiddenWords", &forbiddenWords)
	FindReturnDecoded(filter, "whiteList", &whiteList)
	intervalGetUpdate := 3
	// log.Println("=014257=", whiteList)

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

func isAdmin(message Message, whiteList []WhiteList) bool {
	user := message.From.Username
	adminsMap := make(map[string]bool)
	for _, item := range whiteList {
		if item.Type == "admin" {
			adminsMap[item.Content] = true
		}
	}
	return adminsMap[user]
}

func GetToApi(route string) (io.ReadCloser, error) {
	base := "https://api.telegram.org/bot" + os.Getenv("token") + "/" + route
	res, err := http.Get(base)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return nil, err
	}
	return res.Body, nil
}

func handleBotCommand(entities []Entities, messageText string) {
	switch command := isBotCommands(entities, messageText); command {
	case "/settings":
		categoryCommands := []InlineKeyboardButton{
			{Text: "Запрещенные слова", CallbackData: "forbidden_words"},
			{Text: "WhiteList", CallbackData: "whitelist"},
		}
		sendKeybordMessage(2135753546, "Выберите категорию:", [][]InlineKeyboardButton{categoryCommands})

	case "/status":
		fmt.Println(command, "is great!")
	default:

	}

}
func editKeybordMessage(text string, chatID int64, messageID int64, keyboard [][]InlineKeyboardButton) {
	message := SendMessage{
		ChatID:    chatID,
		Text:      text,
		MessageID: messageID,
		ReplyMarkup: &InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		},
	}
	messageJSON, _ := json.Marshal(message)
	if _, err := PostToApi("editMessageReplyMarkup", messageJSON); err != nil {
		log.Println("=52a1d9=", err)
	}
}
func sendKeybordMessage(chatID int64, text string, keyboard [][]InlineKeyboardButton) {
	message := SendMessage{
		ChatID: chatID,
		Text:   text,
		ReplyMarkup: &InlineKeyboardMarkup{
			InlineKeyboard: keyboard,
		},
	}

	messageJSON, _ := json.Marshal(message)
	if _, err := PostToApi("sendMessage", messageJSON); err != nil {
		log.Println("=52a1d9=", err)
	}
}
func isBotCommands(entities []Entities, messageText string) string {
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
func handleEntities(entities []Entities, messageText string, whiteList []WhiteList) bool {

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
func isContainsForbiddenWord(message string, forbiddenWords []ForbiddenWords) bool {
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

func getUpdate(offset int) (GetUpdates, error) {
	resBody, err := GetToApi(fmt.Sprintf("getUpdates?offset=%d", offset))
	if err != nil {
		return GetUpdates{}, fmt.Errorf("error fetching data: %s", err)
	}
	defer resBody.Close()

	var response GetUpdates
	if err := json.NewDecoder(resBody).Decode(&response); err != nil {
		return GetUpdates{}, fmt.Errorf("error decoding JSON: %s", err)
	}

	return response, nil
}

func Start(host string) {
	http.ListenAndServe(host, nil)
}

func PostToApi(route string, requestBody []byte) (io.ReadCloser, error) {
	base := "https://api.telegram.org/bot" + os.Getenv("token") + "/" + route
	res, err := http.Post(base, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return nil, err
	}
	return res.Body, nil
}

func sendMenu() error {
	// Формирование меню с описанием
	commands := []SetMenu{
		{"/status", "Статус"},
		{"/settings", "Настройки"},
		{"/history", "История"},
	}
	commandsJSON := map[string]interface{}{
		"commands": commands,
	}

	// Формирование сообщения с меню
	requestBody, err := json.Marshal(commandsJSON)
	if err != nil {
		return err
	}
	if res, err := PostToApi("setMyCommands", requestBody); err != nil {
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
