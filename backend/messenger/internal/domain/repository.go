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
	Persist(tx *sql.Tx, user *User) error
	GetByID(tx *sql.Tx, id string) (*User, error)
	GetByEmail(tx *sql.Tx, email string) (*User, error)
	GetByIDAndRefreshToken(tx *sql.Tx, id, token string) (*User, error)
	GetByIDAndAccessToken(tx *sql.Tx, id, token string) (*User, error)
	GetCount(*sql.Tx) (int, error)
	GetByLimitAndOffsetExceptUserID(tx *sql.Tx, userID string, limit, offset int) ([]*User, error)
	GetByPrefixOfNameAndSurname(tx *sql.Tx, prefix string) ([]*User, error)
	UpdateByID(tx *sql.Tx, user *User) error
	CompareError(err error, number uint16) bool
}

type MessengerRepository interface {
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
	DoesUserExist(ctx context.Context, userID string) (bool, error)
	GetAllUsers(ctx context.Context) ([]string, error)
	Persist(ctx context.Context, userID string) error
}

type WSPoolRepository interface {
	AddConnection(userID string, conn net.Conn)
	RemoveConnection(userID string, conn net.Conn)
}
