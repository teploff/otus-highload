package implementation

import (
	"auth/internal/config"
	"auth/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type authService struct {
	repository  domain.UserRepository
	jwtSettings struct {
		secret                 []byte
		accessTokenExpire      time.Duration
		refreshTokenTimeExpire time.Duration
	}
}

func NewAuthService(rep domain.UserRepository, cfg config.JWTConfig) *authService {
	return &authService{
		repository: rep,
		jwtSettings: struct {
			secret                 []byte
			accessTokenExpire      time.Duration
			refreshTokenTimeExpire time.Duration
		}{
			secret:                 []byte(cfg.Secret),
			accessTokenExpire:      cfg.AccessTokenTimeExpire,
			refreshTokenTimeExpire: cfg.RefreshTokenTimeExpire},
	}
}

func (a *authService) SignUp(ctx context.Context, profile *domain.User) error {
	tx, err := a.repository.GetTx(ctx)
	if err != nil {
		return err
	}

	if err = a.repository.Persist(tx, profile); err != nil {
		if a.repository.CompareError(err, domain.DuplicateKeyErrNumber) {
			return fmt.Errorf(fmt.Sprintf("user with email: %s already exist", profile.Email))
		}

		return err
	}

	return a.repository.CommitTx(tx)
}

func (a *authService) SignIn(ctx context.Context, credentials *domain.Credentials) (*domain.TokenPair, error) {
	tx, err := a.repository.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	user, err := a.repository.GetByEmail(tx, credentials.Email)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, fmt.Errorf("incorrect username or password")
	case errors.Is(err, sql.ErrNoRows):
		return nil, err

	}
	if err != nil {
		return nil, err
	}

	if !user.DoesPasswordMatch(credentials.Password) {
		return nil, fmt.Errorf("incorrect username or password")
	}

	tokenPair, err := a.createTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("creating jwt token's pair error, %w", err)
	}

	user.AccessToken = &tokenPair.AccessToken
	user.RefreshToken = &tokenPair.RefreshToken

	if err = a.repository.UpdateByID(tx, user); err != nil {
		return nil, err
	}

	return &tokenPair, a.repository.CommitTx(tx)
}

func (a *authService) createTokenPair(user *domain.User) (domain.TokenPair, error) {
	var tokenPair domain.TokenPair

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(a.jwtSettings.accessTokenExpire).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"iss": "auth_service",
		"aud": user.ID,
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(a.jwtSettings.refreshTokenTimeExpire).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"iss": "auth_service",
		"aud": user.ID,
	})

	var err error
	tokenPair.AccessToken, err = accessToken.SignedString(a.jwtSettings.secret)

	if err != nil {
		return tokenPair, fmt.Errorf("creating jwt token's pair error: %w", err)
	}

	tokenPair.RefreshToken, err = refreshToken.SignedString(a.jwtSettings.secret)
	if err != nil {
		return tokenPair, fmt.Errorf("creating jwt token's pair error: %w", err)
	}

	return tokenPair, nil
}

func (a *authService) parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return a.jwtSettings.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing token err: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("parsing token err")
}

func (a *authService) RefreshToken(ctx context.Context, token string) (*domain.TokenPair, error) {
	tx, err := a.repository.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	claims, err := a.parseToken(token)
	if err != nil {
		return nil, err
	}

	userID := claims["aud"].(string)
	user, err := a.repository.GetByIDAndRefreshToken(tx, userID, token)
	if err != nil {
		return nil, err
	}

	newTokenPair, err := a.createTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("unknown server error: %w", err)
	}

	user.AccessToken = &newTokenPair.AccessToken
	user.RefreshToken = &newTokenPair.RefreshToken

	if err = a.repository.UpdateByID(tx, user); err != nil {
		return nil, err
	}

	return &newTokenPair, a.repository.CommitTx(tx)
}

func (a *authService) Authenticate(ctx context.Context, token string) (*domain.User, error) {
	tx, err := a.repository.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	claims, err := a.parseToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	userID := claims["aud"].(string)
	user, err := a.repository.GetByIDAndAccessToken(tx, userID, token)
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	return user, a.repository.CommitTx(tx)
}

func (a *authService) GetUserIDByEmail(ctx context.Context, email string) (string, error) {
	tx, err := a.repository.GetTx(ctx)
	if err != nil {
		return "", err
	}

	user, err := a.repository.GetByEmail(tx, email)
	if err != nil {
		return "", err
	}

	return user.ID, a.repository.CommitTx(tx)
}

func (a *authService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	tx, err := a.repository.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	user, err := a.repository.GetByID(tx, id)
	if err != nil {
		return nil, err
	}

	return user, a.repository.CommitTx(tx)
}

func (a *authService) SearchByAnthroponym(ctx context.Context, anthroponym, userID string, limit, offset int) ([]*domain.User, int, error) {
	tx, err := a.repository.GetTx(ctx)
	if err != nil {
		return nil, 0, err
	}

	users, count, err := a.repository.GetByAnthroponym(tx, anthroponym, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return users, count, a.repository.CommitTx(tx)
}

func (a *authService) GetUsersByIDs(ctx context.Context, ids []string) ([]*domain.User, error) {
	tx, err := a.repository.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	users, err := a.repository.GetByIDs(tx, ids)
	if err != nil {
		return nil, err
	}

	return users, a.repository.CommitTx(tx)
}
