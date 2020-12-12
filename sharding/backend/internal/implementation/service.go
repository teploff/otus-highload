package implementation

import (
	"backend/internal/config"
	"backend/internal/domain"
	"context"
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
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

func (a *authService) Authenticate(ctx context.Context, token string) (string, error) {
	tx, err := a.repository.GetTx(ctx)
	if err != nil {
		return "", err
	}

	claims, err := a.parseToken(token)
	if err != nil {
		return "", fmt.Errorf("invalid token")
	}

	userID := claims["aud"].(string)
	_, err = a.repository.GetByIDAndAccessToken(tx, userID, token)
	if err != nil {
		return "", fmt.Errorf("invalid token")
	}

	return userID, a.repository.CommitTx(tx)
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

type socialService struct {
	repository domain.UserRepository
}

func NewSocialService(rep domain.UserRepository) *socialService {
	return &socialService{
		repository: rep,
	}
}

func (s *socialService) GetQuestionnaires(ctx context.Context, userID string, limit, offset int) ([]*domain.Questionnaire, int, error) {
	tx, err := s.repository.GetTx(ctx)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repository.GetCount(tx)
	if err != nil {
		return nil, 0, err
	}

	users, err := s.repository.GetByLimitAndOffsetExceptUserID(tx, userID, limit, offset)
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
	return questionnaires, count - 1, s.repository.CommitTx(tx)
}

func (s *socialService) GetQuestionnairesByNameAndSurname(ctx context.Context, prefix string) ([]*domain.Questionnaire, error) {
	tx, err := s.repository.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	users, err := s.repository.GetByPrefixOfNameAndSurname(tx, prefix)
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

	return questionnaires, s.repository.CommitTx(tx)
}

type messengerService struct {
	userRep      domain.UserRepository
	messRep      domain.MessengerRepository
	cacheRep     domain.CacheRepository
	shardingCfg  config.ShardingConfig
	userActivity *senderActivity
}

func NewMessengerService(userRep domain.UserRepository, messengerRep domain.MessengerRepository,
	cacheRep domain.CacheRepository, cfg config.ShardingConfig) *messengerService {
	return &messengerService{
		userRep:      userRep,
		messRep:      messengerRep,
		cacheRep:     cacheRep,
		shardingCfg:  cfg,
		userActivity: newSenderActivity(cfg.MaxMsgFreq),
	}
}

func (m *messengerService) CreateChat(ctx context.Context, masterID, slaveID string) (string, error) {
	tx, err := m.userRep.GetTx(ctx)
	if err != nil {
		return "", err
	}

	_, err = m.userRep.GetByID(tx, slaveID)
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

	return chatID, m.userRep.CommitTx(tx)
}

func (m *messengerService) GetChat(_ context.Context, masterID, slaveID string) (*domain.Chat, error) {
	chat, err := m.messRep.GetChatWithCompanion(masterID, slaveID)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (m *messengerService) GetChats(_ context.Context, userID string, limit, offset int) ([]*domain.Chat, int, error) {
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
	var shardID int

	isLadyGaga, err := m.cacheRep.DoesUserExist(ctx, userID)
	if err != nil {
		return err
	}

	if !isLadyGaga {
		if isLadyGaga = m.userActivity.DoesUserLadyGaga(userID, len(messages)); isLadyGaga {
			if err = m.cacheRep.Persist(ctx, userID); err != nil {
				return err
			}
		}
	}

	if isLadyGaga {
		shardID = m.shardingCfg.LadyGagaShardID
	} else {
		shardID = int(binary.BigEndian.Uint64(uuid.NewV4().Bytes()) % uint64(m.shardingCfg.CountNodes))
	}

	return m.messRep.SendMessages(shardID, userID, chatID, messages)
}

func (m *messengerService) GetMessages(_ context.Context, userID, chatID string, limit, offset int) ([]*domain.Message, int, error) {
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
