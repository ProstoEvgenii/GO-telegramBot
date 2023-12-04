package main

import (
	"GO-chatModeratorTg/db"
	"GO-chatModeratorTg/models"
	"GO-chatModeratorTg/server"
	updates "GO-chatModeratorTg/updates_handler"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// host := "127.0.0.1:80"
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Файл .env не найден")
		// host = ":80"
	}
	db.Connect()
	intervalGetUpdate := 3
	intervalGetData := 5

	tickerGetUpdates := time.NewTicker(time.Duration(intervalGetUpdate) * time.Second)
	tickerGetData := time.NewTicker(time.Duration(intervalGetData) * time.Minute)
	defer tickerGetUpdates.Stop()
	offset := 668578288
	for {
		select {
		case <-tickerGetUpdates.C:
			log.Println("=Получаю обновления каждые=", intervalGetUpdate, "Секунды")

			response, err := getUpdate(offset)
			if err != nil {
				log.Println("=038abf=", err)
			}

			for _, item := range response.Result {
				updates.UpdatesHandler(item)
				offset = item.UpdateID + 1

			}
		}
		select {
		case <-tickerGetData.C:
			log.Println("=Получаю данные из базы каждые=", intervalGetUpdate, "Секунды")
			updates.GetWhiteListAndForbiddeWords()
		}
	}
}

func Start(host string) {
	http.ListenAndServe(host, nil)
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
