package http

import "time"

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
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type QuestionnairesRequest struct {
	Limit  *int `json:"limit" binding:"required"`
	Offset int  `json:"offset"`
}

type AuthorizationHeader struct {
	AccessToken string `json:"access_token" binding:"required" header:"Authorization"`
}
