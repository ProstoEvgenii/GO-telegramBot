package main

import (
	"GO-chatModeratorTg/db"
	"GO-chatModeratorTg/models"
	"GO-chatModeratorTg/server"
	updates "GO-chatModeratorTg/updates_handler"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	host := "127.0.0.1:80"
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Файл .env не найден")
		host = ":80"
	}
	db.Connect()
	intervalGetUpdate := 1
	intervalGetData := 2
	updates.GetWhiteListAndForbiddeWords()
	offset := 668578288
	go runTickers(intervalGetUpdate, intervalGetData, offset)
	server.Start(host)

}
func runTickers(intervalGetUpdate, intervalGetData, offset int) {
	tickerGetUpdates := time.NewTicker(time.Duration(intervalGetUpdate) * time.Second)
	tickerGetData := time.NewTicker(time.Duration(intervalGetData) * time.Second)
	defer tickerGetUpdates.Stop()
	for {
		select {
		case <-tickerGetUpdates.C:
			response, err := getUpdate(offset)
			if err != nil {
				log.Println("=Ошибка получения Update=", err)
			}
			for _, item := range response.Result {
				updates.UpdatesHandler(item)
				offset = item.UpdateID + 1

			}
		case <-tickerGetData.C:
			updates.GetWhiteListAndForbiddeWords()
		}

	}

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
