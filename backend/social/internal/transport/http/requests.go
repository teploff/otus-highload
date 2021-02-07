package http

type SearchProfileByAnthroponymRequest struct {
	Anthroponym string `json:"anthroponym" form:"anthroponym" binding:"required"`
	Limit       *int   `json:"limit" form:"limit"`
	Offset      *int   `json:"offset" form:"offset"`
}

type FriendshipRequest struct {
	FriendsID []string `json:"friends_id" binding:"required"`
}

type GetNewsRequest struct {
	Limit  *int `json:"limit" form:"limit"`
	Offset *int `json:"offset" form:"offset"`
}

type CreateNewsRequest struct {
	News []string `json:"news" binding:"required"`
}

type AuthorizationHeader struct {
	AccessToken string `json:"access_token" binding:"required" header:"Authorization"`
}

type WSRequest struct {
	AccessToken string `json:"token"  binding:"required" form:"token"`
}

type GetUserFriendships struct {
	UserID string `json:"user_id" binding:"required" form:"user_id"`
}

type GetUserIDByAccessTokenRequest struct {
	Token string `json:"token"`
}

type GetUsersByAnthroponymRequest struct {
	Token       string `json:"token"`
	Anthroponym string `json:"anthroponym"`
	Offset      int    `json:"offset"`
	Limit       int    `json:"limit"`
}

type GetUsersByIDsRequest struct {
	Token   string   `json:"token"`
	UserIDs []string `json:"user_ids"`
}
