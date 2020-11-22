package domain

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewCredentials(email string, password string) (*Credentials, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}
	return &Credentials{Email: email, Password: string(hash)}, nil
}

func (c *Credentials) DoesPasswordMatch(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password)) == nil
}

type User struct {
	ID string `json:"id"`
	Credentials
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	Birthday     time.Time `json:"birthday"`
	Sex          string    `json:"sex"`
	City         string    `json:"city"`
	Interests    string    `json:"interests"`
	AccessToken  *string   `json:"access_token"`
	RefreshToken *string   `json:"refresh_token"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Questionnaire struct {
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Birthday  time.Time `json:"birthday"`
	Sex       string    `json:"sex"`
	City      string    `json:"city"`
	Interests string    `json:"interests"`
}

type Message struct {
	ID         string    `json:"id"`
	Text       string    `json:"text"`
	Status     int       `json:"status"`
	CreateTime time.Time `json:"create_time"`
	UserID     string    `json:"user_id"`
	ChatID     string    `json:"chat_id"`
}

type Chat struct {
	ID           string         `json:"id"`
	CreateTime   time.Time      `json:"create_time"`
	Participants []*Participant `json:"participants"`
}

type Participant struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}
