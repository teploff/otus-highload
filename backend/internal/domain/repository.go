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
	GetCount(tx *sql.Tx) (int, error)
	GetByLimitAndOffsetExceptUserID(tx *sql.Tx, userID string, limit, offset int) ([]*User, error)
	GetByPrefixOfNameAndSurname(tx *sql.Tx, prefix string) ([]*User, error)
	GetByAnthroponym(tx *sql.Tx, userID, anthroponym string, limit, offset int) ([]*User, int, error)
	UpdateByID(tx *sql.Tx, user *User) error
	CompareError(err error, number uint16) bool
}

type SocialRepository interface {
	GetTx(ctx context.Context) (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
	CreateFriendship(tx *sql.Tx, subjUserID, objUserID string) error
}

type MessengerRepository interface {
	GetTx(ctx context.Context) (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
	CreateChat(tx *sql.Tx, masterID, slaveID string) (string, error)
	GetCountChats(tx *sql.Tx, userID string) (int, error)
	GetChatWithCompanion(tx *sql.Tx, masterID, slaveID string) (*Chat, error)
	GetChatAsParticipant(tx *sql.Tx, userID string) (*Chat, error)
	GetChats(tx *sql.Tx, userID string, limit, offset int) ([]*Chat, error)
	SendMessages(tx *sql.Tx, userID, chatID string, messages []*ShortMessage) error
	GetCountMessages(tx *sql.Tx, chatID string) (int, error)
	GetMessages(tx *sql.Tx, chatID string, limit, offset int) ([]*Message, error)
}

type WSPoolRepository interface {
	AddConnection(userID string, conn net.Conn)
	RemoveConnection(userID string, conn net.Conn)
}
