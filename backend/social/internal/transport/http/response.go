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

type SearchProfileByAnthroponymResponse struct {
	Users []*domain.User `json:"users"`
	Count int            `json:"count"`
}

type GetFriendsResponse struct {
	Friends []*domain.User `json:"friends"`
}

type GetFollowersResponse struct {
	Followers []*domain.User `json:"followers"`
}

type GetUserFriendshipsResponse struct {
	Friendships []*domain.FriendShip
}

type GetNewsResponse struct {
	News  []*domain.News `json:"news"`
	Count int            `json:"count"`
}
