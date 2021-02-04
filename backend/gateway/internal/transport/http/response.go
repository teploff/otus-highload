package http

import (
	"gateway/internal/domain"
	"time"
)

type authenticateResponse struct {
	IsAuthenticated bool `json:"is_authenticated"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CreateChatResponse struct {
	ChatID string `json:"chat_id"`
}

type GetChatResponse struct {
	ID         string    `json:"id"`
	CreateTime time.Time `json:"create_time"`
}

type GetChatsResponse struct {
	Total  int32          `json:"total"`
	Limit  int32          `json:"limit"`
	Offset int32          `json:"offset"`
	Chats  []*domain.Chat `json:"chats"`
}

type GetMessagesResponse struct {
	Total    int32             `json:"total"`
	Limit    int32             `json:"limit"`
	Offset   int32             `json:"offset"`
	Messages []*domain.Message `json:"messages"`
}
