package http

import "auth/internal/domain"

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

type AuthenticateResponse struct {
	IsAuthenticated bool `json:"is_authenticated"`
}

type GetUserByAccessTokenResponse struct {
	*domain.User
}

type GetUsersByAnthroponymResponse struct {
	Count int            `json:"count"`
	Users []*domain.User `json:"users"`
}

type GetUserByIDsResponse struct {
	Users []*domain.User `json:"users"`
}
