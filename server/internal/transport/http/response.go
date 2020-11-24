package http

import (
	"social-network/internal/domain"
	"time"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type EmptyResponse struct {
}

type SignInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type QuestionnairesResponse struct {
	Questionnaires []*domain.Questionnaire `json:"questionnaires"`
	Count          int                     `json:"count"`
}

type CreateChatResponse struct {
	ChatID string `json:"chat_id"`
}

type GetChatResponse struct {
	ID         string    `json:"id"`
	CreateTime time.Time `json:"create_time"`
}

type GetChatsResponse struct {
	Total  int            `json:"total"`
	Limit  *int           `json:"limit"`
	Offset *int           `json:"offset"`
	Chats  []*domain.Chat `json:"chats"`
}
