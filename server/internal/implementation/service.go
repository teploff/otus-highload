package implementation

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"social-network/internal/domain"
	"time"
)

const (
	accessTokenTimeExpire  = time.Hour
	refreshTokenTimeExpire = time.Hour * 24
)

type authService struct {
	repository domain.UserRepository
	secret     []byte
}

func NewAuthService(rep domain.UserRepository, secret string) *authService {
	return &authService{
		repository: rep,
		secret:     []byte(secret),
	}
}

func (a *authService) SignUp(ctx context.Context, profile *domain.User) error {
	return a.repository.Persist(ctx, profile)
}

func (a *authService) SignIn(ctx context.Context, credentials *domain.Credentials) (*domain.TokenPair, error) {
	user, err := a.repository.GetByLogin(ctx, credentials.Login)
	if err != nil {
		return nil, err
	}

	if !user.DoesPasswordMatch(credentials.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	tokenPair, err := a.createTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("creating jwt token's pair error, %w", err)
	}

	user.AccessToken = &tokenPair.AccessToken
	user.RefreshToken = &tokenPair.RefreshToken

	return &tokenPair, a.repository.UpdateByID(ctx, user)
}

func (a *authService) createTokenPair(user *domain.User) (domain.TokenPair, error) {
	var tokenPair domain.TokenPair

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(accessTokenTimeExpire).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"iss": "auth_service",
		"aud": user.ID,
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(refreshTokenTimeExpire).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"iss": "auth_service",
		"aud": user.ID,
	})

	var err error
	tokenPair.AccessToken, err = accessToken.SignedString(a.secret)

	if err != nil {
		return tokenPair, fmt.Errorf("creating jwt token's pair error: %w", err)
	}

	tokenPair.RefreshToken, err = refreshToken.SignedString(a.secret)
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

		return a.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing token err: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("")
}

func (a *authService) RefreshToken(ctx context.Context, token string) (*domain.TokenPair, error) {
	claims, err := a.parseToken(token)
	if err != nil {
		return nil, err
	}

	userID := claims["aud"].(string)
	user, err := a.repository.GetByIDAndRefreshToken(ctx, userID, token)
	if err != nil {
		return nil, err
	}

	newTokenPair, err := a.createTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("unknown server error: %w", err)
	}

	user.AccessToken = &newTokenPair.AccessToken
	user.RefreshToken = &newTokenPair.RefreshToken

	return &newTokenPair, a.repository.UpdateByID(ctx, user)
}

type socialService struct {
	repository domain.UserRepository
}

func NewSocialService(rep domain.UserRepository) *socialService {
	return &socialService{
		repository: rep,
	}
}

func (s *socialService) GetQuestionnaires(ctx context.Context, userID string, limit int) ([]*domain.Questionnaire, int, error) {
	tx, err := s.repository.GetTx(ctx)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repository.GetCount(tx)
	if err != nil {
		return nil, 0, err
	}

	users, err := s.repository.GetByLimitAndOffsetExceptUserID(tx, userID, limit, 0)
	if err != nil {
		return nil, 0, err
	}

	questionnaires := make([]*domain.Questionnaire, 0, len(users))
	for _, user := range users {
		questionnaires = append(questionnaires, &domain.Questionnaire{
			Email:     user.Login,
			Name:      user.Name,
			Surname:   user.Surname,
			Birthday:  user.Birthday,
			Sex:       user.Sex,
			City:      user.City,
			Interests: user.Interests,
		})
	}

	// count - 1: without myself
	return questionnaires, count - 1, tx.Commit()
}
