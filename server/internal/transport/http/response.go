package http

import "social-network/internal/domain"

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
