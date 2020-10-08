package domain

import (
	"context"
	"database/sql"
)

type UserRepository interface {
	GetTx(ctx context.Context) (*sql.Tx, error)
	Persist(ctx context.Context, user *User) error
	GetByLogin(ctx context.Context, login string) (*User, error)
	GetByIDAndRefreshToken(ctx context.Context, id, token string) (*User, error)
	GetCount(tx *sql.Tx) (int, error)
	GetByLimitAndOffsetExceptUserID(tx *sql.Tx, userID string, limit, offset int) ([]*User, error)
	UpdateByID(ctx context.Context, user *User) error
}
