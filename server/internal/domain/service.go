package domain

import (
	"context"
)

type AuthService interface {
	SignUp(ctx context.Context, user *User) error
	SignIn(ctx context.Context, credentials *Credentials) (*TokenPair, error)
	RefreshToken(ctx context.Context, token string) (*TokenPair, error)
	Authenticate(ctx context.Context, token string) (string, error)
}

type SocialService interface {
	GetQuestionnaires(ctx context.Context, userID string, limit, offset int) ([]*Questionnaire, int, error)
	GetQuestionnairesByNameAndSurname(ctx context.Context, prefix string) ([]*Questionnaire, error)
}

type MessengerService interface {
	CreateChat(ctx context.Context, masterID, slaveID string) (string, error)
	GetChats(ctx context.Context, userID string, limit, offset int) ([]*Chat, int, error)
	SendMessages(ctx context.Context, userID, chatID string, messages []*Message) error
	GetMessages(ctx context.Context, userID, chatID string, limit, offset int) ([]*Message, int, error)
}
