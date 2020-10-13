package domain

import (
	"context"
	"database/sql"
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
	UpdateByID(tx *sql.Tx, user *User) error
	CompareError(err error, number uint16) bool
}
