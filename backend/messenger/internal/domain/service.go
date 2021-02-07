package domain

import (
	"context"
	"net"
)

type AuthService interface {
	GetUserByToken(ctx context.Context, token string) (*User, error)
}

type MessengerService interface {
	CreateChat(ctx context.Context, masterToken, slaveID string) (string, error)
	GetChat(ctx context.Context, masterToken, slaveID string) (*Chat, error)
	SendMessages(ctx context.Context, userToken, chatID string, messages []*ShortMessage) error
	GetChats(ctx context.Context, userToken string, offset, limit int) ([]*Chat, int, error)
	GetMessages(ctx context.Context, userToken, chatID string, offset, limit int) ([]*Message, int, error)
}

type WSService interface {
	EstablishConn(ctx context.Context, user *User, coon net.Conn)
	Close()
}
