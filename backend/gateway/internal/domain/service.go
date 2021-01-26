package domain

import "context"

type GRPCMessengerProxyService interface {
	GetChats(offset, limit *int32, ctx context.Context) (*GetChatsResponse, error)
	GetMessages(chatID string, offset, limit *int32, ctx context.Context) (*GetMessagesResponse, error)
}
