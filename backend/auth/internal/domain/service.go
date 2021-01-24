package domain

import (
	"context"
)

type AuthService interface {
	SignUp(ctx context.Context, user *User) error
	SignIn(ctx context.Context, credentials *Credentials) (*TokenPair, error)
	RefreshToken(ctx context.Context, token string) (*TokenPair, error)
	Authenticate(ctx context.Context, token string) (string, error)
	GetUserIDByEmail(ctx context.Context, email string) (string, error)
}
