package domain

import "time"

type Chat struct {
	ID           string         `json:"id"`
	CreateTime   time.Time      `json:"create_time"`
	Participants []*Participant `json:"participants"`
}

type Participant struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type Message struct {
	ID         string    `json:"id"`
	Text       string    `json:"text" binding:"required"`
	Status     string    `json:"status" binding:"oneof=created received read deleted"`
	CreateTime time.Time `json:"create_time"`
	UserID     string    `json:"user_id"`
	ChatID     string    `json:"chat_id"`
}

type GetChatsResponse struct {
	Total  int32   `json:"total"`
	Limit  int32   `json:"limit"`
	Offset int32   `json:"offset"`
	Chats  []*Chat `json:"chats"`
}

type GetMessagesResponse struct {
	Total    int32      `json:"total"`
	Limit    int32      `json:"limit"`
	Offset   int32      `json:"offset"`
	Messages []*Message `json:"messages"`
}
