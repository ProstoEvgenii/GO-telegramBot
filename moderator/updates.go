package moderator

import (
	"GO-chatModeratorTg/db"
	"GO-chatModeratorTg/keyboard"
	"GO-chatModeratorTg/models"
	"GO-chatModeratorTg/server"
	"GO-chatModeratorTg/tg"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var whiteList []models.WhiteList
var forbiddenWords []models.ForbiddenWords
var adminsMap map[string]bool

func init() {
	adminsMap = make(map[string]bool)
}
func RunTickers(intervalGetUpdate, intervalGetData, offset int) {
	tickerGetUpdates := time.NewTicker(time.Duration(intervalGetUpdate) * time.Second)
	tickerGetData := time.NewTicker(time.Duration(intervalGetData) * time.Minute)
	defer tickerGetUpdates.Stop()
	for {
		select {
		case <-tickerGetUpdates.C:
			response, err := getUpdate(offset)
			if err != nil {
				log.Println("=Ошибка получения Update=", err)
			}
			for _, item := range response.Result {
				UpdatesHandler(item)
				offset = item.UpdateID + 1
			}
		case <-tickerGetData.C:
			GetWhiteListAndForbiddeWords()
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

func GetWhiteListAndForbiddeWords() {
	filter := bson.M{}
	db.FindReturnDecoded(filter, "forbiddenWords", &forbiddenWords)
	db.FindReturnDecoded(filter, "whiteList", &whiteList)
	for _, item := range whiteList {
		if item.Type == "admin" {
			adminsMap[item.Content] = true
		}
	}
}

func UpdatesHandler(update models.Result) {
	if update.Callback.ID != "" {
		handleCallback(update.Callback)

	} else {
		handleMessages(update.Message)
	}
}

func handleMessages(message models.Message) {
	if !isAdmin(message.From.Username, whiteList) && message.Chat.Type != "private" {
		//Сообщение пользователя в модерируемый чат.
		handleUserPublicMessages(message)
	}
	if isAdmin(message.From.Username, whiteList) && message.Chat.Type == "private" {
		//Сообщения админа в личку бота.
		handleAdminPrivateMessages(message)
	}
}

func handleCallback(callBack models.CallbackData) {
	if isAdmin(callBack.From.Username, whiteList) && callBack.Message.Chat.Type == "private" {
		keyboard.HandleSettingsCallback(callBack)
	} else if !isAdmin(callBack.From.Username, whiteList) && callBack.Message.Chat.Type == "private" {
		message := models.SendMessage{
			ChatID:    callBack.Message.Chat.ID,
			Text:      "Бот выполняет команды только администратора.",
			ParseMode: "None",
		}
		tg.SendMessage(message)
	}
}

func isAdmin(user string, whiteList []models.WhiteList) bool {
	return adminsMap[user]
}
