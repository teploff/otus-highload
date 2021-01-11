package stan

import "social-network/internal/domain"

// FriendsActionRequest request for AddFriends or DeleteFriends endpoints.
// easyjson:json
type FriendsActionRequest struct {
	Action    string
	UserID    string
	FriendsID []string
}

// NewsPersistRequest request for PublishNews endpoint.
// easyjson:json
type NewsPersistRequest struct {
	OwnerID string         `json:"owner_id"`
	News    []*domain.News `json:"news"`
}
