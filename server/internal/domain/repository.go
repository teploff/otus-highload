package domain

import "context"

type UserRepository interface {
	Persist(ctx context.Context, user *User) error
	GetByLogin(ctx context.Context, login string) (*User, error)
	GetByIDAndRefreshToken(ctx context.Context, id, token string) (*User, error)
	UpdateByID(ctx context.Context, user *User) error
}
