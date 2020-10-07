package domain

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewCredentials(login string, password string) (*Credentials, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}
	return &Credentials{Login: login, Password: string(hash)}, nil
}

func (c *Credentials) DoesPasswordMatch(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password)) == nil
}

type User struct {
	Credentials
	Name         string    `json:"name"`
	Surname      string    `json:"surname"`
	Birthday     time.Time `json:"birthday"`
	Sex          string    `json:"sex"`
	City         string    `json:"city"`
	Interests    string    `json:"interests"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
