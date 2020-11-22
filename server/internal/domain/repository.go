package domain

import (
	"context"
	"database/sql"
	"net"
)

const DuplicateKeyErrNumber = 1062

type UserRepository interface {
	GetTx(ctx context.Context) (*sql.Tx, error)
	Persist(tx *sql.Tx, user *User) error
	GetByEmail(tx *sql.Tx, email string) (*User, error)
	GetByIDAndRefreshToken(tx *sql.Tx, id, token string) (*User, error)
	GetByIDAndAccessToken(tx *sql.Tx, id, token string) (*User, error)
	GetCount(tx *sql.Tx) (int, error)
	GetByLimitAndOffsetExceptUserID(tx *sql.Tx, userID string, limit, offset int) ([]*User, error)
	GetByPrefixOfNameAndSurname(tx *sql.Tx, prefix string) ([]*User, error)
	UpdateByID(tx *sql.Tx, user *User) error
	CompareError(err error, number uint16) bool
}

type MessengerRepository interface {
	GetTx(ctx context.Context) (*sql.Tx, error)
	CreateChat(tx *sql.Tx, masterID, slaveID string) (string, error)
	GetChats(tx *sql.Tx, userID string, limit, offset int) ([]*Chat, int, error)
	SendMessages(tx *sql.Tx, userID, chatID string, messages []*Message) error
	GetMessages(tx *sql.Tx, userID, chatID string, limit, offset int) ([]*Message, int, error)
}

type WSPoolRepository interface {
	AddConnection(userID string, conn net.Conn)
	RemoveConnection(userID string, conn net.Conn)
}
