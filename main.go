package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// host := "127.0.0.1:80"
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Файл .env не найден")
		// host = ":80"
	}
	// Connect()
	words := readForbiddenWords("new.txt")

	interval := 3

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("=Выполняюсь каждые=", interval, "Секунды")
			response, err := getUpdate()
			if err != nil {
				log.Println("=038abf=", err)
			}
			for _, message := range response.Result {
				if containsForbiddenWord(message.Message.Text, words) {
					deleteMessage(message.Message.Chat.ID, message.Message.MessageID)
				}
				if message.Message.Text == "word" {
					log.Printf("=От=%s= %s", message.Message.From.Username, message.Message.Text)
					fmt.Println("=65eca7=", reflect.TypeOf(message.Message.Text))
					// log.Println("=Удалено так как в тексте содержится слово=", word)
					deleteMessage(message.Message.Chat.ID, message.Message.MessageID)
				}
			}
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

func GetToApi(route string) (io.ReadCloser, error) {
	base := "https://api.telegram.org/bot" + os.Getenv("token") + "/" + route
	res, err := http.Get(base)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return nil, err
	}
	return res.Body, nil
}

func containsForbiddenWord(text string, forbiddenWords []string) bool {
	loweredText := strings.Fields(strings.ToLower(text))
	for _, word := range loweredText {
		for _, forbiddenWord := range forbiddenWords {
			if word == forbiddenWord {
				log.Printf("=Сообщение %s будет удалео из-за слова %s", text, forbiddenWord)
				return true
			}
		}
	}
	return false
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

// func containsObscene(text string, obsceneWords []string) bool {
// 	for _, word := range obsceneWords {
// 		// log.Println("=33709f=",word)
// 		if strings.Contains(text, word) {
// 			return true
// 		}
// 	}
// 	return false
// }

// Start(host)
// var chatID int64 = -1002053372425
// bot.Debug = true

// log.Printf("Authorized on account %s", bot.Self.UserName)

// u := tgbotapi.NewUpdate(0)
// u.Timeout = 60

// updates := bot.GetUpdatesChan(u)

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
//
//		}
//	}
func Start(host string) {
	http.ListenAndServe(host, nil)
}
