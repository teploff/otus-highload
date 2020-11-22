package domain

import (
	"context"
	"net"
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

type Messenger interface {
	AddConnection(userID string, conn net.Conn)
	RemoveConnection(userID string, conn net.Conn)
	CreateMessage(msg []byte) error
}
