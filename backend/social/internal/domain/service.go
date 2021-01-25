package domain

import (
	"context"
	"net"
)

type AuthService interface {
	Authenticate(ctx context.Context, token string) (string, error)
}

type ProfileService interface {
	SearchByAnthroponym(ctx context.Context, anthroponym, userID string, limit, offset int) ([]*Questionnaire, int, error)
}

type SocialService interface {
	CreateFriendship(ctx context.Context, userID string, friendsID []string) error
	ConfirmFriendship(ctx context.Context, userID string, friendsID []string) error
	RejectFriendship(ctx context.Context, userID string, friendsID []string) error
	BreakFriendship(ctx context.Context, userID string, friendsID []string) error
	GetFriends(ctx context.Context, userID string) ([]*Questionnaire, error)
	GetFollowers(ctx context.Context, userID string) ([]*Questionnaire, error)
	RetrieveNews(ctx context.Context, userID string, limit, offset int) ([]*News, int, error)
	PublishNews(ctx context.Context, userID string, newsContent []string) error
}

type CacheService interface {
	AddFriends(ctx context.Context, userID string, friendsID []string) error
	DeleteFriends(ctx context.Context, userID string, friendsID []string) error
	AddNews(ctx context.Context, userID string, news []*News) error
}

type WSService interface {
	EstablishConn(ctx context.Context, userID string, coon net.Conn)
	SendNews(ctx context.Context, ownerID string, news []*News) error
	Close()
}
