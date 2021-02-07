package domain

import (
	"context"
	"database/sql"
	"net"
)

type SocialRepository interface {
	GetTx(ctx context.Context) (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
	CreateFriendship(tx *sql.Tx, userID string, friendsID []string) error
	ConfirmFriendship(tx *sql.Tx, userID string, friendsID []string) error
	RejectFriendship(tx *sql.Tx, userID string, friendsID []string) error
	BreakFriendship(tx *sql.Tx, userID string, friendsID []string) error
	GetFriends(tx *sql.Tx, userID string) ([]string, error)
	GetFollowers(tx *sql.Tx, userID string) ([]string, error)
	GetUserFriendships(tx *sql.Tx, userID string) ([]*FriendShip, error)
	GetNews(tx *sql.Tx, friends []*User, offset, limit int) ([]*News, int, error)
	PublishNews(tx *sql.Tx, userID string, news []*News) error
}

type SocialCacheRepository interface {
	PersistFriend(ctx context.Context, userID string, friendsID []string) error
	DeleteFriend(ctx context.Context, userID string, friendsID []string) error
	RetrieveFriendsID(ctx context.Context, userID string) ([]string, error)
	PersistNews(ctx context.Context, userID string, news []*News) error
	RetrieveNews(ctx context.Context, userID string) ([]*News, error)
}

type WSPoolRepository interface {
	AddConnection(userID string, conn net.Conn)
	RemoveConnection(userID string, conn net.Conn)
	RetrieveConnByUserID(userID string) []net.Conn
	FlushConnections()
}
