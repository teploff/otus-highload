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

type ProfileService interface {
	SearchByAnthroponym(ctx context.Context, anthroponym, userID string, limit, offset int) ([]*Questionnaire, int, error)
}

type SocialService interface {
	CreateFriendship(ctx context.Context, userID, friendID string) error
	ConfirmFriendship(ctx context.Context, userID, friendID string) error
	RejectFriendship(ctx context.Context, userID, friendID string) error
	BreakFriendship(ctx context.Context, userID, friendID string) error
	GetFriends(ctx context.Context, userID string) ([]*Questionnaire, error)
	GetFollowers(ctx context.Context, userID string) ([]*Questionnaire, error)
	GetQuestionnaires(ctx context.Context, userID string, limit, offset int) ([]*Questionnaire, int, error)
	GetQuestionnairesByNameAndSurname(ctx context.Context, prefix string) ([]*Questionnaire, error)
}

type MessengerService interface {
	CreateChat(ctx context.Context, masterID, slaveID string) (string, error)
	GetChat(ctx context.Context, masterID, slaveID string) (*Chat, error)
	SendMessages(ctx context.Context, userID, chatID string, messages []*ShortMessage) error
	GetChats(ctx context.Context, userID string, limit, offset int) ([]*Chat, int, error)
	GetMessages(ctx context.Context, userID, chatID string, limit, offset int) ([]*Message, int, error)
}
