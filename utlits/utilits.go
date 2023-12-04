package utlits

import (
	"GO-chatModeratorTg/models"
	"GO-chatModeratorTg/server"
	"bufio"
	"encoding/json"
	"log"
	"os"
)

func sendMenu() error {
	// Формирование меню с описанием
	commands := []models.SetMenu{
		{Command: "/status", Description: "Статус"},
		{Command: "/settings", Description: "Настройки"},
		{Command: "/history", Description: "История"},
	}
	commandsJSON := map[string]interface{}{
		"commands": commands,
	}

	requestBody, err := json.Marshal(commandsJSON)
	if err != nil {
		return err
	}
	if res, err := server.PostToApi("setMyCommands", requestBody); err != nil {
		log.Println("=c77107=", res)
		return err
	}
	return nil
}
func readFile(filename string) []string {
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
