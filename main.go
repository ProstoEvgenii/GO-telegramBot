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

	var forbiddenWords []ForbiddenWords
	var whiteList []WhiteList
	FindReturnDecoded(filter, "forbiddenWords", &forbiddenWords)
	FindReturnDecoded(filter, "whiteList", &whiteList)
	intervalGetUpdate := 3
	log.Println("=014257=", whiteList)

	ticker := time.NewTicker(time.Duration(intervalGetUpdate) * time.Second)
	defer ticker.Stop()
	offset := 668578288
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

				if !isAdmin(item.Message, whiteList) {
					loweredText := strings.ToLower(item.Message.Text)

					if len(item.Message.Entities) != 0 {
						if !handleEntities(item.Message.Entities, loweredText, whiteList) {
							deleteMessage(item.Message.Chat.ID, item.Message.MessageID)
						}

					} else if isContainsForbiddenWord(loweredText, forbiddenWords) {
						deleteMessage(item.Message.Chat.ID, item.Message.MessageID)
					}
				}
			}
		}
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
