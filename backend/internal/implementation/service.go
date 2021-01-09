package implementation

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"social-network/internal/config"
	"social-network/internal/domain"
	"time"
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

type profileService struct {
	repository domain.UserRepository
}

func NewProfileService(repository domain.UserRepository) *profileService {
	return &profileService{repository: repository}
}

func (p *profileService) SearchByAnthroponym(ctx context.Context, anthroponym, userID string, limit, offset int) ([]*domain.Questionnaire, int, error) {
	tx, err := p.repository.GetTx(ctx)
	if err != nil {
		return nil, 0, err
	}

	users, count, err := p.repository.GetByAnthroponym(tx, anthroponym, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	questionnaires := make([]*domain.Questionnaire, 0, len(users))
	for _, user := range users {
		questionnaires = append(questionnaires, &domain.Questionnaire{
			ID:               user.ID,
			Email:            user.Email,
			Name:             user.Name,
			Surname:          user.Surname,
			Birthday:         user.Birthday,
			Sex:              user.Sex,
			City:             user.City,
			Interests:        user.Interests,
			FriendshipStatus: user.FriendshipStatus,
		})
	}

	return questionnaires, count, p.repository.CommitTx(tx)
}

type socialService struct {
	userRepository   domain.UserRepository
	socialRepository domain.SocialRepository
}

func NewSocialService(userRep domain.UserRepository, socialRep domain.SocialRepository) *socialService {
	return &socialService{
		userRepository:   userRep,
		socialRepository: socialRep,
	}
}

func (s *socialService) CreateFriendship(ctx context.Context, userID, friendID string) error {
	tx, err := s.socialRepository.GetTx(ctx)
	if err != nil {
		return err
	}

	if err = s.socialRepository.CreateFriendship(tx, userID, friendID); err != nil {
		return err
	}

	return s.socialRepository.CommitTx(tx)
}

func (s *socialService) ConfirmFriendship(ctx context.Context, userID, friendID string) error {
	tx, err := s.socialRepository.GetTx(ctx)
	if err != nil {
		return err
	}

	if err = s.socialRepository.CreateFriendship(tx, userID, friendID); err != nil {
		return err
	}

	return s.socialRepository.CommitTx(tx)
}

func (s *socialService) RejectFriendship(ctx context.Context, userID, friendID string) error {
	tx, err := s.socialRepository.GetTx(ctx)
	if err != nil {
		return err
	}

	if err = s.socialRepository.CreateFriendship(tx, userID, friendID); err != nil {
		return err
	}

	return s.socialRepository.CommitTx(tx)
}

func (s *socialService) BreakFriendship(ctx context.Context, userID, friendID string) error {
	tx, err := s.socialRepository.GetTx(ctx)
	if err != nil {
		return err
	}

	if err = s.socialRepository.BreakFriendship(tx, userID, friendID); err != nil {
		return err
	}

	return s.socialRepository.CommitTx(tx)
}

func (s *socialService) GetFriends(ctx context.Context, userID string) ([]*domain.Questionnaire, error) {
	tx, err := s.socialRepository.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	users, err := s.socialRepository.GetFriends(tx, userID)
	if err != nil {
		return nil, err
	}

	questionnaires := make([]*domain.Questionnaire, 0, len(users))
	for _, user := range users {
		questionnaires = append(questionnaires, &domain.Questionnaire{
			ID:               user.ID,
			Email:            user.Email,
			Name:             user.Name,
			Surname:          user.Surname,
			Birthday:         user.Birthday,
			Sex:              user.Sex,
			City:             user.City,
			Interests:        user.Interests,
			FriendshipStatus: user.FriendshipStatus,
		})
	}

	return questionnaires, s.socialRepository.CommitTx(tx)
}

func (s *socialService) GetFollowers(ctx context.Context, userID string) ([]*domain.Questionnaire, error) {
	tx, err := s.socialRepository.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	users, err := s.socialRepository.GetFollowers(tx, userID)
	if err != nil {
		return nil, err
	}

	questionnaires := make([]*domain.Questionnaire, 0, len(users))
	for _, user := range users {
		questionnaires = append(questionnaires, &domain.Questionnaire{
			ID:               user.ID,
			Email:            user.Email,
			Name:             user.Name,
			Surname:          user.Surname,
			Birthday:         user.Birthday,
			Sex:              user.Sex,
			City:             user.City,
			Interests:        user.Interests,
			FriendshipStatus: user.FriendshipStatus,
		})
	}

	return questionnaires, s.socialRepository.CommitTx(tx)
}

func (s *socialService) GetQuestionnaires(ctx context.Context, userID string, limit, offset int) ([]*domain.Questionnaire, int, error) {
	tx, err := s.userRepository.GetTx(ctx)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.userRepository.GetCount(tx)
	if err != nil {
		return nil, 0, err
	}

	users, err := s.userRepository.GetByLimitAndOffsetExceptUserID(tx, userID, limit, offset)
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
	return questionnaires, count - 1, s.userRepository.CommitTx(tx)
}

func (s *socialService) GetQuestionnairesByNameAndSurname(ctx context.Context, prefix string) ([]*domain.Questionnaire, error) {
	tx, err := s.userRepository.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	users, err := s.userRepository.GetByPrefixOfNameAndSurname(tx, prefix)
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

	return questionnaires, s.userRepository.CommitTx(tx)
}

type messengerService struct {
	userRep domain.UserRepository
	messRep domain.MessengerRepository
}

func NewMessengerService(userRep domain.UserRepository, messengerRep domain.MessengerRepository) *messengerService {
	return &messengerService{
		userRep: userRep,
		messRep: messengerRep,
	}
}

func (m *messengerService) CreateChat(ctx context.Context, masterID, slaveID string) (string, error) {
	tx, err := m.messRep.GetTx(ctx)
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

	chatID, err := m.messRep.CreateChat(tx, masterID, slaveID)
	if err != nil {
		return "", err
	}

	return chatID, m.messRep.CommitTx(tx)
}

func (m *messengerService) GetChat(ctx context.Context, masterID, slaveID string) (*domain.Chat, error) {
	tx, err := m.messRep.GetTx(ctx)
	if err != nil {
		return nil, err
	}

	chat, err := m.messRep.GetChatWithCompanion(tx, masterID, slaveID)
	if err != nil {
		return nil, err
	}

	return chat, m.messRep.CommitTx(tx)
}

func (m *messengerService) GetChats(ctx context.Context, userID string, limit, offset int) ([]*domain.Chat, int, error) {
	tx, err := m.messRep.GetTx(ctx)
	if err != nil {
		return nil, 0, err
	}

	total, err := m.messRep.GetCountChats(tx, userID)
	if err != nil {
		return nil, 0, err
	}

	chats, err := m.messRep.GetChats(tx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return chats, total, m.messRep.CommitTx(tx)
}

func (m *messengerService) SendMessages(ctx context.Context, userID, chatID string, messages []*domain.ShortMessage) error {
	tx, err := m.messRep.GetTx(ctx)
	if err != nil {
		return err
	}

	err = m.messRep.SendMessages(tx, userID, chatID, messages)
	if err != nil {
		return err
	}

	return m.messRep.CommitTx(tx)
}

func (m *messengerService) GetMessages(ctx context.Context, userID, chatID string, limit, offset int) ([]*domain.Message, int, error) {
	tx, err := m.messRep.GetTx(ctx)
	if err != nil {
		return nil, 0, err
	}

	_, err = m.messRep.GetChatAsParticipant(tx, userID)
	if err != nil {
		return nil, 0, err
	}

	total, err := m.messRep.GetCountMessages(tx, chatID)
	if err != nil {
		return nil, 0, err
	}

	messages, err := m.messRep.GetMessages(tx, chatID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return messages, total, m.messRep.CommitTx(tx)
}
