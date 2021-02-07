package domain

import (
	"context"
	"net"
)

type AuthService interface {
	GetUserByToken(ctx context.Context, token string) (*User, error)
	GetUsersByAnthroponym(ctx context.Context, token, anthroponym string, offset, limit int) ([]*User, int, error)
	GetUsersByIDs(ctx context.Context, ids []string) ([]*User, error)
}

type ProfileService interface {
	SearchByAnthroponym(ctx context.Context, token, anthroponym string, offset, limit int) ([]*User, int, error)
}

type SocialService interface {
	CreateFriendship(ctx context.Context, userID string, friendsID []string) error
	ConfirmFriendship(ctx context.Context, userID string, friendsID []string) error
	RejectFriendship(ctx context.Context, userID string, friendsID []string) error
	BreakFriendship(ctx context.Context, userID string, friendsID []string) error
	GetFriends(ctx context.Context, userID string) ([]*User, error)
	GetFollowers(ctx context.Context, userID string) ([]*User, error)
	RetrieveNews(ctx context.Context, userID string, offset, limit int) ([]*News, int, error)
	PublishNews(ctx context.Context, userID string, newsContent []string) error
}

type CacheService interface {
	AddFriends(ctx context.Context, userID string, friendsID []string) error
	DeleteFriends(ctx context.Context, userID string, friendsID []string) error
	AddNews(ctx context.Context, userID string, news []*News) error
}

type WSService interface {
	EstablishConn(ctx context.Context, user *User, coon net.Conn)
	SendNews(ctx context.Context, ownerID string, news []*News) error
	Close()
}
