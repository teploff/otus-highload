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

type FriendShip struct {
	MasterUserID string `json:"master_user_id"`
	SlaveUserID  string `json:"slave_user_id"`
	Status       string `json:"status"`
}

// easyjson:json
type News struct {
	ID    string `json:"-"`
	Owner struct {
		Name    string `json:"name"`
		Surname string `json:"surname"`
		Sex     string `json:"sex"`
	} `json:"owner"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"create_time"`
}
