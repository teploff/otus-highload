package domain

import (
	"context"
	"database/sql"
	"net"
)

const DuplicateKeyErrNumber = 1062

type UserRepository interface {
	GetTx(ctx context.Context) (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
	GetByID(tx *sql.Tx, id string) (*User, error)
	GetByAnthroponym(tx *sql.Tx, userID, anthroponym string, limit, offset int) ([]*User, int, error)
	CompareError(err error, number uint16) bool
}

type SocialRepository interface {
	GetTx(ctx context.Context) (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
	CreateFriendship(tx *sql.Tx, userID string, friendsID []string) error
	ConfirmFriendship(tx *sql.Tx, userID string, friendsID []string) error
	RejectFriendship(tx *sql.Tx, userID string, friendsID []string) error
	BreakFriendship(tx *sql.Tx, userID string, friendsID []string) error
	GetFriends(tx *sql.Tx, userID string) ([]*User, error)
	GetFollowers(tx *sql.Tx, userID string) ([]*User, error)
	GetNews(tx *sql.Tx, userID string, limit, offset int) ([]*News, int, error)
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
