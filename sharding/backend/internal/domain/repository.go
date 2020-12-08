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
	Persist(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByIDAndRefreshToken(id, token string) (*User, error)
	GetByIDAndAccessToken(id, token string) (*User, error)
	GetCount() (int, error)
	GetByLimitAndOffsetExceptUserID(userID string, limit, offset int) ([]*User, error)
	GetByPrefixOfNameAndSurname(prefix string) ([]*User, error)
	UpdateByID(user *User) error
	CompareError(err error, number uint16) bool
}

type MessengerRepository interface {
	GetTx(ctx context.Context) (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
	CreateChat(masterID, slaveID string) (string, error)
	GetCountChats(userID string) (int, error)
	GetChatWithCompanion(masterID, slaveID string) (*Chat, error)
	GetChatAsParticipant(userID string) (*Chat, error)
	GetChats(userID string, limit, offset int) ([]*Chat, error)
	SendMessages(shardID int, userID, chatID string, messages []*ShortMessage) error
	GetCountMessages(chatID string) (int, error)
	GetMessages(chatID string, limit, offset int) ([]*Message, error)
}

type CacheRepository interface {
	GetLadyGagaUsers(ctx context.Context) ([]string, error)
	Persist(ctx context.Context, userID string) error
}

type WSPoolRepository interface {
	AddConnection(userID string, conn net.Conn)
	RemoveConnection(userID string, conn net.Conn)
}
