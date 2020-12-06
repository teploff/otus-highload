package http

import (
	"backend/internal/domain"
	"time"
)

type SignUpRequest struct {
	Email     string    `json:"email" binding:"required"`
	Password  string    `json:"password" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	Surname   string    `json:"surname" binding:"required"`
	Birthday  time.Time `json:"birthday" binding:"required"`
	Sex       string    `json:"sex" binding:"required"`
	City      string    `json:"city" binding:"required"`
	Interests string    `json:"interests" binding:"required"`
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type GetAllQuestionnairesRequest struct {
	Limit  int  `json:"limit" binding:"required"`
	Offset *int `json:"offset" binding:"required"`
}

type GetQuestionnairesByNameAndSurnameRequest struct {
	Prefix string `json:"prefix" form:"prefix" binding:"required"`
	Limit  int    `json:"limit" form:"limit" binding:"required"`
	Offset *int   `json:"offset" form:"offset" binding:"required"`
}

type AuthorizationHeader struct {
	AccessToken string `json:"access_token" binding:"required" header:"Authorization"`
}

type CreateChatRequest struct {
	CompanionID string `json:"companion_id" binding:"required"`
}

type GetChatRequest struct {
	CompanionID string `json:"companion_id" form:"companion_id" binding:"required"`
}

type GetChatsRequest struct {
	Limit  *int `json:"limit" form:"limit"`
	Offset *int `json:"offset" form:"offset"`
}

type SendMessagesRequest struct {
	ChatID   string                 `json:"chat_id" binding:"required"`
	Messages []*domain.ShortMessage `json:"messages" binding:"required,dive"`
}

type GetMessagesRequest struct {
	ChatID string `json:"chat_id" form:"chat_id" binding:"required"`
	Limit  *int   `json:"limit" form:"limit"`
	Offset *int   `json:"offset" form:"offset"`
}
