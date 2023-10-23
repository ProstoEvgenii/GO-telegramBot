package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6385727068:AAFN2Jtlu4LOEoMCm04S9dS4I_9T2EHBQ4M")
	if err != nil {
		log.Panic(err)
	}
	var chatID int64 = -1002053372425
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	obsceneWords := readForbiddenWords("new.txt")
	// log.Println("=2763d5=", obsceneWords)
	for update := range updates {
		if update.Message != nil { // If we got a message_
			// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			log.Println("=15021a=", update.Message.From.UserName)

			// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			// msg.ReplyToMessageID = update.Message.MessageID
			// log.Println("=64411e=", update.Message.Text, update.Message.Chat.ID)

			// bot.Send(msg)
			messageText := strings.ToLower(update.Message.Text)
			if containsObscene(messageText, obsceneWords) {
				// msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("ты нехороший человек %s", update.Message.From.FirstName))
				msg2 := tgbotapi.NewMessage(update.Message.From.ID, fmt.Sprintf("ты нехороший человек %s", update.Message.From.FirstName))
				// bot.Send(msg)
				bot.Send(msg2)
				deleteMessage(bot, chatID, update.Message.MessageID)
			}
			// if update.Message.Text == "Привет" && update.Message.Chat.ID == chatID {
			// 	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("ты нехороший человек %s", update.Message.From.UserName))
			// 	bot.Send(msg)
			// 	deleteMessage(bot, chatID, update.Message.MessageID)

			// }
		}
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
func containsObscene(text string, obsceneWords []string) bool {
	for _, word := range obsceneWords {
		// log.Println("=33709f=",word)
		if strings.Contains(text, word) {
			return true
		}
	}
	return false
}
func deleteMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	msg := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error deleting message: %v", err)
	}
}

// func deleteMessage(chatID int64, messageID int) {
// 	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/deleteMessage", botToken)

// 	message := map[string]interface{}{
// 		"chat_id":    chatID,
// 		"message_id": messageID,
// 	}

// 	body, err := json.Marshal(message)
// 	if err != nil {
// 		log.Println("Ошибка при создании JSON-запроса для удаления сообщения:", err)
// 		return
// 	}

// 	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
// 	if err != nil {
// 		log.Println("Ошибка при удалении сообщения:", err)
// 	}
// 	defer resp.Body.Close()
// }
