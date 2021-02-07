package domain

import (
	"context"
)

type AuthService interface {
	SignUp(ctx context.Context, user *User) error
	SignIn(ctx context.Context, credentials *Credentials) (*TokenPair, error)
	RefreshToken(ctx context.Context, token string) (*TokenPair, error)
	Authenticate(ctx context.Context, token string) (*User, error)
	GetUserIDByEmail(ctx context.Context, email string) (string, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	SearchByAnthroponym(ctx context.Context, anthroponym, userID string, limit, offset int) ([]*User, int, error)
	GetUsersByIDs(ctx context.Context, ids []string) ([]*User, error)
}
