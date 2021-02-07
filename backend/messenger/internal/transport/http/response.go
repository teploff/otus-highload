package http

import "messenger/internal/domain"

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type GetUserIDByAccessTokenResponse struct {
	*domain.User
}

type EmptyResponse struct {
}
