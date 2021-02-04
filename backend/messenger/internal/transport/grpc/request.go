package grpc

type CreateChatRequest struct {
	MasterToken string `json:"master_token"`
	SlaveID     string `json:"slave_id"`
}

type GetChatRequest struct {
	UserToken   string `json:"user_token"`
	CompanionID string `json:"companion_id"`
}

type GetChatsRequest struct {
	UserToken string `json:"user_token"`
	Limit     *int   `json:"limit"`
	Offset    *int   `json:"offset"`
}

type GetMessagesRequest struct {
	UserToken string `json:"user_token"`
	ChatID    string `json:"chat_id"`
	Limit     *int   `json:"limit"`
	Offset    *int   `json:"offset"`
}
