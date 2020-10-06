package implementation

import (
	"context"
	"social-network/internal/domain"
)

type authService struct {
	repository domain.ProfileRepository
}

func NewAuthService(rep domain.ProfileRepository) *authService {
	return &authService{
		repository: rep,
	}
}

func (a *authService) SignUp(ctx context.Context, profile *domain.Profile) error {
	return a.repository.Persist(ctx, profile)
}

func (a *authService) SignIn(ctx context.Context, credentials *domain.Credentials) (*domain.TokenPair, error) {
	panic("implement me")
}

func (a *authService) RefreshToken(ctx context.Context, token string) (*domain.TokenPair, error) {
	panic("implement me")
}
