package http

import (
	"social/internal/domain"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type EmptyResponse struct {
}

type QuestionnairesResponse struct {
	Questionnaires []*domain.Questionnaire `json:"questionnaires"`
	Count          int                     `json:"count"`
}

type GetFriendsResponse struct {
	Friends []*domain.Questionnaire `json:"friends"`
}

type GetFollowersResponse struct {
	Followers []*domain.Questionnaire `json:"followers"`
}

type GetNewsResponse struct {
	News  []*domain.News `json:"news"`
	Count int            `json:"count"`
}
