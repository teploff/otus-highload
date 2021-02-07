package domain

import "net"

type MessengerRepository interface {
	CreateChat(masterID, slaveID string) (string, error)
	GetCountChats(userID string) (int, error)
	GetChatWithCompanion(masterID, slaveID string) (*Chat, error)
	GetChatAsParticipant(userID string) (*Chat, error)
	GetParticipantsByChatID(userID, chatID string) ([]string, error)
	GetChats(userID string, limit, offset int) ([]*Chat, error)
	PersistMessages(userID, chatID string, messages []*ShortMessage) error
	GetCountMessages(chatID string) (int, error)
	GetMessages(chatID string, limit, offset int) ([]*Message, error)
}

type WSPoolRepository interface {
	AddConnection(userID string, conn net.Conn)
	RemoveConnection(userID string, conn net.Conn)
	RetrieveConnByUserID(userID string) []net.Conn
	FlushConnections()
}
