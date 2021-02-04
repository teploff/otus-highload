package domain

import (
	"context"
	"net"
)

type AuthService interface {
	Authenticate(ctx context.Context, token string) (*User, error)
}

type MessengerService interface {
	CreateChat(ctx context.Context, masterToken, slaveID string) (string, error)
	GetChat(ctx context.Context, masterToken, slaveID string) (*Chat, error)
	SendMessages(ctx context.Context, userToken, chatID string, messages []*ShortMessage) error
	GetChats(ctx context.Context, userToken string, limit, offset int) ([]*Chat, int, error)
	GetMessages(ctx context.Context, userToken, chatID string, limit, offset int) ([]*Message, int, error)
}

type WSService interface {
	EstablishConn(ctx context.Context, userID string, coon net.Conn)
	Close()
}
