package domain

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

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

type Message struct {
	ID         string    `json:"id"`
	Text       string    `json:"text" binding:"required"`
	Status     string    `json:"status" binding:"oneof=created received read deleted"`
	CreateTime time.Time `json:"create_time"`
	UserID     string    `json:"user_id"`
	ChatID     string    `json:"chat_id"`
}

type GetChatsResponse struct {
	Total  int32   `json:"total"`
	Limit  int32   `json:"limit"`
	Offset int32   `json:"offset"`
	Chats  []*Chat `json:"chats"`
}

type GetMessagesResponse struct {
	Total    int32      `json:"total"`
	Limit    int32      `json:"limit"`
	Offset   int32      `json:"offset"`
	Messages []*Message `json:"messages"`
}

type ServerAvailableList struct {
	servers []string

	sync.Mutex
}

func NewServerAvailableList() *ServerAvailableList {
	return &ServerAvailableList{servers: make([]string, 0, 1)}
}

func (s *ServerAvailableList) GetAddr() (string, error) {
	s.Lock()
	defer s.Unlock()

	if len(s.servers) == 0 {
		return "", fmt.Errorf("all severs are not available")
	}

	return s.servers[rand.Intn(len(s.servers))], nil
}

func (s *ServerAvailableList) Update(servers []string) {
	s.Lock()
	defer s.Unlock()

	s.servers = servers
}
