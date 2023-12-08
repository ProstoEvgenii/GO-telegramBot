package moderator

import (
	"GO-chatModeratorTg/models"
	"GO-chatModeratorTg/tg"
	"log"
	"regexp"
	"strings"
)

func handleUserPublicMessages(message models.Message) {
	loweredText := strings.ToLower(message.Text)
	if len(message.Entities) != 0 {
		if !handleEntities(message.Entities, loweredText, whiteList) {
			//Проверка сущностей используемых в сообщении(ссылки, упоминания)
			tg.DeleteMessage(message.Chat.ID, message.MessageID)
		}
	} else if isContainsForbiddenWord(loweredText, forbiddenWords) {
		//Проверка на запрещенные слова
		tg.DeleteMessage(message.Chat.ID, message.MessageID)
	}
}

// Бот проверяет сущности чата(ссылки, упоминания) для дальнейшей модерации.
func handleEntities(entities []models.Entities, messageText string, whiteList []models.WhiteList) bool {
	for _, entity := range entities {
		if entity.Type == "mention" {
			//Обрабатывать разрешенные упоминания
			start := entity.Offset
			stop := start + entity.Length
			mention := messageText[start+1 : stop]
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

// Бот разбивает сообщение на слайс строк и сверяет каждую строку с содержимым карты запрещенных слов.
func isContainsForbiddenWord(message string, forbiddenWords []models.ForbiddenWords) bool {
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
