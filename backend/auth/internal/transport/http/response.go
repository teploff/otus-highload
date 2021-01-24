package http

import (
	"auth/internal/domain"
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

type GetUserIDByEmailResponse struct {
	UserID string `json:"user_id"`
}

type QuestionnairesResponse struct {
	Questionnaires []*domain.Questionnaire `json:"questionnaires"`
	Count          int                     `json:"count"`
}
