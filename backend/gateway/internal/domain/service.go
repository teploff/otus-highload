package domain

import "context"

type GRPCMessengerProxyService interface {
	CreateChat(ctx context.Context, userToken, companionID string) (string, error)
	GetChats(ctx context.Context, userToken string, offset, limit *int32) (*GetChatsResponse, error)
	GetMessages(ctx context.Context, userToken, chatID string, offset, limit *int32) (*GetMessagesResponse, error)
}
