package http

type WSRequest struct {
	AccessToken string `json:"token"  binding:"required" form:"token"`
}

type AuthorizationHeader struct {
	AccessToken string `json:"access_token" binding:"required" header:"Authorization"`
}

type CreateChatRequest struct {
	CompanionID string `json:"companion_id" binding:"required"`
}

type GetChatRequest struct {
	CompanionID string `json:"companion_id" form:"companion_id" binding:"required"`
}

type GetChatsRequest struct {
	Limit  *int32 `json:"limit" form:"limit"`
	Offset *int32 `json:"offset" form:"offset"`
}

type GetMessagesRequest struct {
	ChatID string `json:"chat_id" form:"chat_id" binding:"required"`
	Limit  *int32 `json:"limit" form:"limit"`
	Offset *int32 `json:"offset" form:"offset"`
}
