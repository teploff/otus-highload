package stan

// FriendsActionRequest request for AddFriends or DeleteFriends endpoints.
// easyjson:json
type FriendsActionRequest struct {
	Action    string
	UserID    string
	FriendsID []string
}

// NewsActionRequest request for PublishNews endpoint.
// easyjson:json
type NewsActionRequest struct {
}
