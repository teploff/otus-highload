package grpc

import (
	"messenger/internal/domain"
	"time"
)

type CreateChatResponse struct {
	ChatID string `json:"chat_id"`
}

type GetChatResponse struct {
	ID         string    `json:"id"`
	CreateTime time.Time `json:"create_time"`
}

type GetChatsResponse struct {
	Total  int            `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
	Chats  []*domain.Chat `json:"chats"`
}

type GetMessagesResponse struct {
	Total    int               `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
	Messages []*domain.Message `json:"messages"`
}
