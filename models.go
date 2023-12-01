package main

type GetUpdates struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID int     `json:"update_id"`
		Message  Message `json:"message"`
	} `json:"result"`
}
type Message struct {
	MessageID int `json:"message_id"`
	From      struct {
		ID           int    `json:"id"`
		IsBot        bool   `json:"is_bot"`
		FirstName    string `json:"first_name"`
		Username     string `json:"username"`
		LanguageCode string `json:"language_code"`
	} `json:"from"`
	Chat struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		Username  string `json:"username"`
		Type      string `json:"type"`
	} `json:"chat"`
	Date     int        `json:"date"`
	Text     string     `json:"text"`
	Entities []Entities `json:"entities"`
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
	Name string `bson:"word"`
}

// Если нужно декодировать result - создаю структуру и указываю ее в Result как тип данных
type ChatMemberResult struct {
	Status string `json:"status"`
}
