package implementation

import (
	"backend/internal/config"
	"backend/internal/domain"
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
	if err := a.repository.Persist(profile); err != nil {
		if a.repository.CompareError(err, domain.DuplicateKeyErrNumber) {
			return fmt.Errorf(fmt.Sprintf("user with email: %s already exist", profile.Email))
		}

		return err
	}

	return nil
}

func (a *authService) SignIn(ctx context.Context, credentials *domain.Credentials) (*domain.TokenPair, error) {
	user, err := a.repository.GetByEmail(credentials.Email)
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

	if err = a.repository.UpdateByID(user); err != nil {
		return nil, err
	}

	return &tokenPair, nil
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
	claims, err := a.parseToken(token)
	if err != nil {
		return nil, err
	}

	userID := claims["aud"].(string)
	user, err := a.repository.GetByIDAndRefreshToken(userID, token)
	if err != nil {
		return nil, err
	}

	newTokenPair, err := a.createTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("unknown server error: %w", err)
	}

	user.AccessToken = &newTokenPair.AccessToken
	user.RefreshToken = &newTokenPair.RefreshToken

	if err = a.repository.UpdateByID(user); err != nil {
		return nil, err
	}

	return &newTokenPair, nil
}

func (a *authService) Authenticate(ctx context.Context, token string) (string, error) {
	claims, err := a.parseToken(token)
	if err != nil {
		return "", fmt.Errorf("invalid token")
	}

	userID := claims["aud"].(string)
	_, err = a.repository.GetByIDAndAccessToken(userID, token)
	if err != nil {
		return "", fmt.Errorf("invalid token")
	}

	return userID, nil
}

func (a *authService) GetUserIDByEmail(ctx context.Context, email string) (string, error) {
	user, err := a.repository.GetByEmail(email)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}

type socialService struct {
	repository domain.UserRepository
}

func NewSocialService(rep domain.UserRepository) *socialService {
	return &socialService{
		repository: rep,
	}
}

func (s *socialService) GetQuestionnaires(ctx context.Context, userID string, limit, offset int) ([]*domain.Questionnaire, int, error) {
	count, err := s.repository.GetCount()
	if err != nil {
		return nil, 0, err
	}

	users, err := s.repository.GetByLimitAndOffsetExceptUserID(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	questionnaires := make([]*domain.Questionnaire, 0, len(users))
	for _, user := range users {
		questionnaires = append(questionnaires, &domain.Questionnaire{
			Email:     user.Email,
			Name:      user.Name,
			Surname:   user.Surname,
			Birthday:  user.Birthday,
			Sex:       user.Sex,
			City:      user.City,
			Interests: user.Interests,
		})
	}

	// count - 1: without myself
	return questionnaires, count - 1, nil
}

func (s *socialService) GetQuestionnairesByNameAndSurname(ctx context.Context, prefix string) ([]*domain.Questionnaire, error) {
	users, err := s.repository.GetByPrefixOfNameAndSurname(prefix)
	if err != nil {
		return nil, err
	}

	questionnaires := make([]*domain.Questionnaire, 0, len(users))
	for _, user := range users {
		questionnaires = append(questionnaires, &domain.Questionnaire{
			Email:     user.Email,
			Name:      user.Name,
			Surname:   user.Surname,
			Birthday:  user.Birthday,
			Sex:       user.Sex,
			City:      user.City,
			Interests: user.Interests,
		})
	}

	return questionnaires, nil
}

type messengerService struct {
	userRep  domain.UserRepository
	messRep  domain.MessengerRepository
	cacheRep domain.CacheRepository
}

func NewMessengerService(userRep domain.UserRepository, messengerRep domain.MessengerRepository,
	cacheRep domain.CacheRepository) *messengerService {
	return &messengerService{
		userRep:  userRep,
		messRep:  messengerRep,
		cacheRep: cacheRep,
	}
}

func (m *messengerService) CreateChat(ctx context.Context, masterID, slaveID string) (string, error) {
	_, err := m.userRep.GetByID(slaveID)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return "", fmt.Errorf("chat companion with id=[%s] doesn't exist", slaveID)
	default:
		return "", err
	}

	chatID, err := m.messRep.CreateChat(masterID, slaveID)
	if err != nil {
		return "", err
	}

	return chatID, nil
}

func (m *messengerService) GetChat(ctx context.Context, masterID, slaveID string) (*domain.Chat, error) {
	chat, err := m.messRep.GetChatWithCompanion(masterID, slaveID)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (m *messengerService) GetChats(ctx context.Context, userID string, limit, offset int) ([]*domain.Chat, int, error) {
	total, err := m.messRep.GetCountChats(userID)
	if err != nil {
		return nil, 0, err
	}

	chats, err := m.messRep.GetChats(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return chats, total, nil
}

func (m *messengerService) SendMessages(ctx context.Context, userID, chatID string, messages []*domain.ShortMessage) error {
	err := m.messRep.SendMessages(0, userID, chatID, messages)
	if err != nil {
		return err
	}

	return nil
}

func (m *messengerService) GetMessages(ctx context.Context, userID, chatID string, limit, offset int) ([]*domain.Message, int, error) {
	_, err := m.messRep.GetChatAsParticipant(userID)
	if err != nil {
		return nil, 0, err
	}

	total, err := m.messRep.GetCountMessages(chatID)
	if err != nil {
		return nil, 0, err
	}

	messages, err := m.messRep.GetMessages(chatID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}
