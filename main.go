package main

import (
	"GO-chatModeratorTg/db"
	"GO-chatModeratorTg/moderator"
	"GO-chatModeratorTg/server"
	"GO-chatModeratorTg/tg"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	host := "127.0.0.1:80"
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Файл .env не найден")
		host = ":80"
	}
	db.Connect()
	tg.SendMenu()
	intervalGetUpdate := 1
	intervalGetData := 15
	moderator.GetWhiteListAndForbiddeWords()
	offset := 668578288
	go moderator.RunTickers(intervalGetUpdate, intervalGetData, offset)
	server.Start(host)

}
