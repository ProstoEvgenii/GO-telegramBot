package main

import "time"

type GetUpdates struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID int          `json:"update_id"`
		Message  Message      `json:"message"`
		Callback CallbackData `json:"callback_query"`
	} `json:"result"`
}

type Message struct {
	MessageID int64 `json:"message_id"`
	From      struct {
		ID           int    `json:"id"`
		IsBot        bool   `json:"is_bot"`
		FirstName    string `json:"first_name"`
		Username     string `json:"username"`
		LanguageCode string `json:"language_code"`
	} `json:"from"`
	Chat     Chat       `json:"chat"`
	Date     int        `json:"date"`
	Text     string     `json:"text"`
	Entities []Entities `json:"entities"`
}

type CallbackData struct {
	ID       string  `json:"id"`
	From     User    `json:"from"`
	Message  Message `json:"message"`
	Data     string  `json:"data"`
	Chat     Chat    `json:"chat"`
	DateTime int     `json:"date"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type Chat struct {
	ID        int64    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}


type Entities struct {
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Type   string `json:"type"`
	URL    string `json:"url"`
}

type GetChatMemberResponse struct {
	Ok     bool             `json:"ok"`
	Result ChatMemberResult `json:"result"`
}
type DeleteMessageResponse struct {
	Ok bool `json:"ok"`
	// Result bool `json:"result"`
}

type ForbiddenWords struct {
	Word string `bson:"word"`
}

type WhiteList struct {
	Type    string    `bson:"type"`
	Content string    `bson:"content"`
	Added   time.Time `bson:"added"`
	AddedBy string    `bson:"addedBy"`
}

// Если нужно декодировать result - создаю структуру и указываю ее в Result как тип данных
type ChatMemberResult struct {
	Status string `json:"status"`
}
