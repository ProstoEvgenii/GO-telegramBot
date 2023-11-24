package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type GetUpdates struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID int `json:"update_id"`
		Message  struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID           int    `json:"id"`
				IsBot        bool   `json:"is_bot"`
				FirstName    string `json:"first_name"`
				Username     string `json:"username"`
				LanguageCode string `json:"language_code"`
			} `json:"from"`
			Chat struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				Username  string `json:"username"`
				Type      string `json:"type"`
			} `json:"chat"`
			Date int    `json:"date"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"result"`
}

type Delete struct {
	Ok     bool `json:"ok"`
	Result bool `json:"result"`
}

func myFunction() {
	// Ваш код здесь
	fmt.Println("Код выполняется каждую минуту")
}
func main() {
	// host := "127.0.0.1:80"
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Файл .env не найден")
		// host = ":80"
	}

	interval := 3

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			myFunction()
			response, err := getUpdate()
			if err != nil {
				log.Println("=038abf=", err)
			}
			for _, message := range response.Result {
				log.Printf("=От=%s= %s", message.Message.From.Username, message.Message.Text)
				if message.Message.Text == "11" {
					deleteMessage(message.Message.Chat.ID, message.Message.MessageID)
				}
			}
		}
	}
}

// Start(host)
// var chatID int64 = -1002053372425
// bot.Debug = true

// log.Printf("Authorized on account %s", bot.Self.UserName)

// u := tgbotapi.NewUpdate(0)
// u.Timeout = 60

// updates := bot.GetUpdatesChan(u)
// obsceneWords := readForbiddenWords("new.txt")
// log.Println("=2763d5=", obsceneWords)
// for update := range updates {
// if update.Message != nil { // If we got a message_
// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
// log.Println("=15021a=", update.Message.From.UserName)

// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
// msg.ReplyToMessageID = update.Message.MessageID
// log.Println("=64411e=", update.Message.Text, update.Message.Chat.ID)

// bot.Send(msg)
// messageText := strings.ToLower(update.Message.Text)
// if containsObscene(messageText, obsceneWords) {
// msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("ты нехороший человек %s", update.Message.From.FirstName))
// msg2 := tgbotapi.NewMessage(update.Message.From.ID, fmt.Sprintf("ты нехороший человек %s", update.Message.From.FirstName))
// bot.Send(msg)
// 	bot.Send(msg2)
// 	deleteMessage(bot, chatID, update.Message.MessageID)
// }
// if update.Message.Text == "Привет" && update.Message.Chat.ID == chatID {
// 	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("ты нехороший человек %s", update.Message.From.UserName))
// 	bot.Send(msg)
// 	deleteMessage(bot, chatID, update.Message.MessageID)

// }
// 		}
// 	}

func Start(host string) {
	http.ListenAndServe(host, nil)
}

func getUpdate() (GetUpdates, error) {
	resBody, err := GetToApi("getUpdates")
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

func GetToApi(route string) (io.ReadCloser, error) {
	base := "https://api.telegram.org/bot" + os.Getenv("token") + "/" + route
	res, err := http.Get(base)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return nil, err
	}
	return res.Body, nil
}

//	func deleteMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
//		msg := tgbotapi.NewDeleteMessage(chatID, messageID)
//		_, err := bot.Send(msg)
//		if err != nil {
//			log.Printf("Error deleting message: %v", err)
//		}
//	}
func deleteMessage(chatID, messageID int) (Delete, error) {
	route := fmt.Sprintf("deleteMessage?chat_id=%d&message_id=%d", chatID, messageID)

	resp, err := GetToApi(route)
	if err != nil {
		return Delete{}, fmt.Errorf("error sending request: %s", err)
	}
	var response Delete

	if err := json.NewDecoder(resp).Decode(&response); err != nil {
		return Delete{}, fmt.Errorf("error decoding JSON: %s", err)
	}
	if !response.Result {
		return Delete{}, fmt.Errorf("error deleting message: %t", response.Result)
	}
	log.Println("=b414b5=", response)
	return response, nil
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
