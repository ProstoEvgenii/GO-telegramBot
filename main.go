package main

import (
	"GO-chatModeratorTg/db"
	"GO-chatModeratorTg/moderator"
	"GO-chatModeratorTg/server"
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
	intervalGetUpdate := 1
	intervalGetData := 15
	offset := 668578288
	moderator.GetWhiteListAndForbiddeWords()
	go moderator.RunTickers(intervalGetUpdate, intervalGetData, offset)
	server.Start(host)

}
