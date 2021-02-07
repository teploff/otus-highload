package implementation

import (
	"database/sql"
	"encoding/binary"
	"messenger/internal/config"
	"messenger/internal/domain"
	"messenger/internal/infrastructure/clickhouse"
	"net"
	"sync"
	"time"

	querybuilder "github.com/rtsoftSG/plugin/toolbox/database"
	uuid "github.com/satori/go.uuid"
)

type messengerRepository struct {
	conn *clickhouse.Storage
	cfg  config.ShardingConfig
}

func NewMessengerRepository(conn *clickhouse.Storage, cfg config.ShardingConfig) *messengerRepository {
	return &messengerRepository{
		conn: conn,
		cfg:  cfg,
	}
}

func (m *messengerRepository) CreateChat(masterID, slaveID string) (string, error) {
	var chatID string

	err := m.conn.DB().QueryRow(`
		SELECT
			id
		FROM
			chat
		WHERE hasAll(participants, [toUUID(?), toUUID(?)])`, masterID, slaveID).Scan(&chatID)

	switch err {
	case nil:
		return chatID, nil
	case sql.ErrNoRows:
		chatID = uuid.NewV4().String()
	default:
		return "", err
	}

	now := time.Now().UTC()
	cmd := querybuilder.NewInsertCommand("chat").
		WithFields(
			querybuilder.NewField("datetime", now),
			querybuilder.NewField("create_time", now.UnixNano()),
			querybuilder.NewField("id", chatID),
			querybuilder.NewField("participants", []string{masterID, slaveID}),
		)
	m.conn.Insert(cmd)

	return chatID, nil
}

func (m *messengerRepository) GetCountChats(userID string) (int, error) {
	var count int

	err := m.conn.DB().QueryRow(`
		SELECT 
			count(*)
		FROM chat
		WHERE hasAll(participants, [toUUID(?)])`, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *messengerRepository) GetChatWithCompanion(masterID, slaveID string) (*domain.Chat, error) {
	//var chat domain.Chat
	//
	//err := m.conn.QueryRow(`
	//	SELECT
	//		C1.id, C1.create_time
	//	FROM (
	//		SELECT
	//			chat.id, chat.create_time
	//		FROM user
	//		JOIN user_chat
	//			ON user.id = user_chat.user_id
	//		JOIN chat
	//			ON user_chat.chat_id = chat.id
	//		WHERE user.id = ?
	//	) as C1
	//	JOIN user_chat AS UC1
	//		ON C1.id = UC1.chat_id
	//	where UC1.user_id = ?`, masterID, slaveID).Scan(&chat.ID, &chat.CreateTime)
	//if err != nil {
	//	return nil, err
	//}

	return nil, nil
}

func (m *messengerRepository) GetParticipantsByChatID(userID, chatID string) ([]string, error) {
	participants := make([]string, 0, 1)

	err := m.conn.DB().QueryRow(`
		SELECT
			participants
		FROM
			chat
		WHERE id = ? AND hasAll(participants, [toUUID(?)])`, chatID, userID).Scan(&participants)
	if err != nil {
		return nil, err
	}

	for index, id := range participants {
		if id == userID {
			participants[index] = participants[len(participants)-1]
			participants[len(participants)-1] = ""
			participants = participants[:len(participants)-1]
		}
	}

	return participants, nil
}

func (m *messengerRepository) GetChatAsParticipant(userID string) (*domain.Chat, error) {
	var (
		chat     domain.Chat
		unixNano int64
	)

	err := m.conn.DB().QueryRow(`
		SELECT
			id, create_time
		FROM
			chat
		WHERE hasAll(participants, [toUUID(?)])`, userID).Scan(&chat.ID, &unixNano)
	if err != nil {
		return nil, err
	}

	chat.CreateTime = time.Unix(0, unixNano)

	return &chat, nil
}

func (m *messengerRepository) GetChats(userID string, limit, offset int) ([]*domain.Chat, error) {
	chats := make([]*domain.Chat, 0, 10)

	//rows, err := m.conn.Query(`
	//	SELECT
	//		chat.id, chat.create_time
	//	FROM user
	//	JOIN user_chat
	//		ON user.id = user_chat.user_id
	//	JOIN chat
	//		ON user_chat.chat_id = chat.id
	//	WHERE user.id = ?
	//	ORDER BY chat.id
	//	LIMIT ? OFFSET ?`, userID, limit, offset)
	//if err != nil {
	//	return nil, err
	//}
	//
	//type chatRow struct {
	//	id         string
	//	createTime time.Time
	//}
	//chatRows := make([]*chatRow, 0)
	//
	//for rows.Next() {
	//	var row chatRow
	//	if err = rows.Scan(&row.id, &row.createTime); err != nil {
	//		return nil, err
	//	}
	//
	//	chatRows = append(chatRows, &row)
	//}
	//rows.Close()
	//
	//for _, chat := range chatRows {
	//	rows, err = m.conn.Query(`
	//	SELECT
	//		user.id, user.name, user.surname
	//	FROM user_chat
	//	JOIN user
	//		ON user_chat.user_id = user.id
	//	WHERE user_chat.chat_id = ? AND user_chat.user_id != ?`, chat.id, userID)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	for rows.Next() {
	//		var user struct {
	//			ID      string
	//			Name    string
	//			Surname string
	//		}
	//		if err = rows.Scan(&user.ID, &user.Name, &user.Surname); err != nil {
	//			return nil, err
	//		}
	//
	//		exist := false
	//		for _, c := range chats {
	//			if c.ID == chat.id {
	//				c.Participants = append(c.Participants, &domain.Participant{
	//					ID:      user.ID,
	//					Name:    user.Name,
	//					Surname: user.Surname,
	//				})
	//
	//				exist = true
	//			}
	//		}
	//
	//		if !exist {
	//			c := &domain.Chat{
	//				ID:           chat.id,
	//				CreateTime:   chat.createTime,
	//				Participants: make([]*domain.Participant, 0),
	//			}
	//			c.Participants = append(c.Participants, &domain.Participant{
	//				ID:      user.ID,
	//				Name:    user.Name,
	//				Surname: user.Surname,
	//			})
	//
	//			chats = append(chats, c)
	//		}
	//	}
	//	rows.Close()
	//}

	return chats, nil
}

func (m *messengerRepository) PersistMessages(userID, chatID string, messages []*domain.ShortMessage) error {
	shardID := int(binary.BigEndian.Uint64([]byte(userID)) % uint64(m.cfg.CountNodes))

	for _, msg := range messages {
		now := time.Now().UTC()

		cmd := querybuilder.NewInsertCommand("message").
			WithFields(
				querybuilder.NewField("datetime", now),
				querybuilder.NewField("create_time", now.UnixNano()),
				querybuilder.NewField("id", uuid.NewV4().String()),
				querybuilder.NewField("text", msg.Text),
				querybuilder.NewField("status", "created"),
				querybuilder.NewField("user_id", userID),
				querybuilder.NewField("chat_id", chatID),
				querybuilder.NewField("shard_id", shardID),
			)
		m.conn.Insert(cmd)
	}

	return nil
}

func (m *messengerRepository) GetCountMessages(chatID string) (int, error) {
	var count int
	//select count(*) from (select max(create_time) from message where chat_id = "188a4d72-64a7-4dd4-beae-df8ca11fce70" group by id) AS a;
	err := m.conn.DB().QueryRow(`
		SELECT 
			count(*)
		FROM (
		    SELECT
		    	MAX(create_time)
		    FROM message
		    WHERE chat_id = ?
		    GROUP BY id
		) AS MSG`, chatID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *messengerRepository) GetMessages(chatID string, limit, offset int) ([]*domain.Message, error) {
	messages := make([]*domain.Message, 0, 10)

	rows, err := m.conn.DB().Query(`
		SELECT id, text, status, user_id, max(create_time) as create_time FROM message WHERE chat_id = ? GROUP BY id, text, status, user_id LIMIT ?,?`, chatID, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			message  domain.Message
			unixNano int64
		)

		if err = rows.Scan(&message.ID, &message.Text, &message.Status, &message.UserID, &unixNano); err != nil {
			return nil, err
		}

		message.CreateTime = time.Unix(0, unixNano)
		messages = append(messages, &message)
	}

	return messages, nil
}

type wsPoolRepository struct {
	conns map[string][]net.Conn

	m sync.Mutex
}

func NewWSPoolRepository() *wsPoolRepository {
	return &wsPoolRepository{
		conns: make(map[string][]net.Conn),
		m:     sync.Mutex{},
	}
}

func (w *wsPoolRepository) AddConnection(userID string, conn net.Conn) {
	w.m.Lock()
	defer w.m.Unlock()

	w.conns[userID] = append(w.conns[userID], conn)
}

func (w *wsPoolRepository) RemoveConnection(userID string, conn net.Conn) {
	w.m.Lock()
	defer w.m.Unlock()

	for index, connection := range w.conns[userID] {
		if connection == conn {
			w.conns[userID][index] = w.conns[userID][len(w.conns[userID])-1]
			w.conns[userID][len(w.conns[userID])-1] = nil
			w.conns[userID] = w.conns[userID][:len(w.conns[userID])-1]

			return
		}
	}
}

func (w *wsPoolRepository) FlushConnections() {
	w.m.Lock()
	defer w.m.Unlock()

	for userID, userConns := range w.conns {
		for _, conn := range userConns {
			conn.Close()
		}
		w.conns[userID] = nil
	}
}

func (w *wsPoolRepository) RetrieveConnByUserID(userID string) []net.Conn {
	w.m.Lock()
	defer w.m.Unlock()

	if _, ok := w.conns[userID]; ok {
		return w.conns[userID]
	}

	return make([]net.Conn, 0, 1)
}
