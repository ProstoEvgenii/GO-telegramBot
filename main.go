package main

import (
	"bufio"
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
	offset := 668578288
	for {
		select {
		case <-ticker.C:
			log.Println("=Выполняюсь каждые=", interval, "Секунды")

			response, err := getUpdate(offset)
			if err != nil {
				log.Println("=038abf=", err)
			}

			for _, item := range response.Result {
				offset = item.UpdateID + 1
				log.Println("=d25e31=", item.Message)
				if containsForbiddenWord(item.Message, words) {
					if item.Message.From.Username != "dmitriibelov" {
						deleteMessage(item.Message.Chat.ID, item.Message.MessageID)
					}
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

func containsForbiddenWord(message Message, forbiddenWords []string) bool {
	loweredText := strings.Fields(strings.ToLower(message.Text))
	regClean := regexp.MustCompile("[^a-zA-Zа-яА-Я0-9]+")
	regUsername := regexp.MustCompile(`@([a-zA-Z0-9]+)`)
	for _, messageWord := range loweredText {
		for _, forbiddenWord := range forbiddenWords {
			loweredForbiddenWord := strings.ToLower(forbiddenWord)
			cleanedMessageWord := regClean.ReplaceAllString(messageWord, "")
			if cleanedMessageWord == loweredForbiddenWord {
				// log.Println("=7c7444=", cleanedMessageWord)
				log.Printf("=Сообщение будет удалено из-за слова %s", forbiddenWord)
				return true
			} else if regUsername.MatchString(messageWord) {
				log.Printf("=Сообщение  %s удалено. \nИспользование логина", messageWord)
				return true
			}
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
