package domain

import (
	"context"
	"net"
)

type AuthService interface {
	Authenticate(ctx context.Context, token string) (string, error)
}

type MessengerService interface {
	CreateChat(ctx context.Context, masterID, slaveID string) (string, error)
	GetChat(ctx context.Context, masterID, slaveID string) (*Chat, error)
	SendMessages(ctx context.Context, userID, chatID string, messages []*ShortMessage) error
	GetChats(ctx context.Context, userID string, limit, offset int) ([]*Chat, int, error)
	GetMessages(ctx context.Context, userID, chatID string, limit, offset int) ([]*Message, int, error)
}

type WSService interface {
	EstablishConn(ctx context.Context, userID string, coon net.Conn)
	//SendNews(ctx context.Context, ownerID string, news []*Message) error
	Close()
}
