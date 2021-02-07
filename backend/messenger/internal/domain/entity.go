package domain

import (
	"time"
)

type User struct {
	ID               string    `json:"id"`
	Email            string    `json:"email"`
	Name             string    `json:"name"`
	Surname          string    `json:"surname"`
	Birthday         time.Time `json:"birthday"`
	Sex              string    `json:"sex"`
	City             string    `json:"city"`
	Interests        string    `json:"interests"`
	FriendshipStatus string    `json:"friendship_status"`
}

type ShortMessage struct {
	Text   string `json:"text" binding:"required"`
	Status string `json:"status" binding:"oneof=created received read deleted"`
}

type Message struct {
	ID string `json:"id"`
	ShortMessage
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
