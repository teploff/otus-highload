package implementation

import (
	"backend/internal/domain"
	"backend/internal/infrastructure/clickhouse"
	wstransport "backend/internal/transport/ws"
	"context"
	"database/sql"
	"net"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	querybuilder "github.com/rtsoftSG/plugin/toolbox/database"
	uuid "github.com/satori/go.uuid"
)

type userRepository struct {
	conn *sql.DB
}

func NewUserRepository(conn *sql.DB) *userRepository {
	return &userRepository{conn: conn}
}

func (p *userRepository) GetTx(ctx context.Context) (*sql.Tx, error) {
	return p.conn.BeginTx(ctx, nil)
}

func (p *userRepository) CommitTx(tx *sql.Tx) error {
	return tx.Commit()
}

func (p *userRepository) Persist(tx *sql.Tx, user *domain.User) error {
	_, err := tx.Exec(`
		INSERT 
			INTO user (email, password, name, surname, birthday, sex, city, interests) 
		VALUES
			( ?, ?, ?, ?, ?, ?, ?, ?)`, user.Email, user.Password, user.Name, user.Surname, user.Birthday, user.Sex,
		user.City, user.Interests)
	if err != nil {
		tx.Rollback()

		return err
	}

	return nil
}

func (p *userRepository) GetByID(tx *sql.Tx, id string) (*domain.User, error) {
	var user domain.User

	err := tx.QueryRow(`
		SELECT
			id, email, password, name, surname, sex, birthday, city, interests, access_token, refresh_token
		FROM
			 user 
		WHERE 
			  id = ?`, id).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Surname,
		&user.Sex, &user.Birthday, &user.City, &user.Interests, &user.AccessToken, &user.RefreshToken)
	if err != nil {
		tx.Rollback()

		return nil, err
	}

	return &user, nil
}

func (p *userRepository) GetByEmail(tx *sql.Tx, email string) (*domain.User, error) {
	var user domain.User

	err := tx.QueryRow(`
		SELECT
			id, email, password, name, surname, sex, birthday, city, interests, access_token, refresh_token
		FROM
			 user 
		WHERE 
			  email = ?`, email).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Surname,
		&user.Sex, &user.Birthday, &user.City, &user.Interests, &user.AccessToken, &user.RefreshToken)
	if err != nil {
		tx.Rollback()

		return nil, err
	}

	return &user, nil
}

func (p *userRepository) GetByIDAndRefreshToken(tx *sql.Tx, id, token string) (*domain.User, error) {
	var user domain.User

	err := tx.QueryRow(`
		SELECT
			id, email, password, name, surname, sex, birthday, city, interests, access_token, refresh_token
		FROM
			user
		WHERE
			id = ? AND refresh_token = ?`, id, token).Scan(&user.ID, &user.Email, &user.Password, &user.Name,
		&user.Surname, &user.Sex, &user.Birthday, &user.City, &user.Interests, &user.AccessToken,
		&user.RefreshToken)
	if err != nil {
		tx.Rollback()

		return nil, err
	}

	return &user, nil
}

func (p *userRepository) GetByIDAndAccessToken(tx *sql.Tx, id, token string) (*domain.User, error) {
	var user domain.User

	err := tx.QueryRow(`
		SELECT
			id, email, password, name, surname, sex, birthday, city, interests, access_token, refresh_token
		FROM
			user
		WHERE
			id = ? AND access_token = ?`, id, token).Scan(&user.ID, &user.Email, &user.Password, &user.Name,
		&user.Surname, &user.Sex, &user.Birthday, &user.City, &user.Interests, &user.AccessToken,
		&user.RefreshToken)
	if err != nil {
		tx.Rollback()

		return nil, err
	}

	return &user, nil
}

func (p *userRepository) UpdateByID(tx *sql.Tx, user *domain.User) error {
	_, err := tx.Exec(`
		UPDATE 
			user
		SET
		    email = ?, password = ?, name = ?, surname = ?, birthday = ?, sex = ?, city = ?, interests = ?,
		    access_token = ?, refresh_token = ?, update_time = ?
		WHERE
		    id = ?`, user.Email, user.Password, user.Name, user.Surname, user.Birthday, user.Sex,
		user.City, user.Interests, user.AccessToken, user.RefreshToken, time.Now().UTC(), user.ID)
	if err != nil {
		tx.Rollback()

		return err
	}

	return nil
}

func (p *userRepository) GetCount(tx *sql.Tx) (int, error) {
	var count int

	if err := tx.QueryRow(`SELECT count(*) FROM user`).Scan(&count); err != nil {
		tx.Rollback()

		return 0, err
	}

	return count, nil
}

func (p *userRepository) GetByLimitAndOffsetExceptUserID(tx *sql.Tx, userID string, limit, offset int) ([]*domain.User, error) {
	users := make([]*domain.User, 0, 10)

	rows, err := tx.Query(`
		SELECT
			id, email, password, name, surname, sex, birthday, city, interests, access_token, refresh_token
		FROM
		    user
		WHERE 
			  id != ?
		ORDER BY create_time
		LIMIT ? OFFSET ?`, userID, limit, offset)
	if err != nil {
		tx.Rollback()

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := new(domain.User)

		if err = rows.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Surname, &user.Sex, &user.Birthday,
			&user.City, &user.Interests, &user.AccessToken, &user.RefreshToken); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (p *userRepository) GetByPrefixOfNameAndSurname(tx *sql.Tx, prefix string) ([]*domain.User, error) {
	users := make([]*domain.User, 0, 100)

	rows, err := tx.Query(`
		SELECT
			id, email, password, name, surname, sex, birthday, city, interests, access_token, refresh_token
		FROM
		    user
		WHERE 
			  name LIKE ? AND surname LIKE ?
		ORDER BY id`, prefix+"%", prefix+"%")
	if err != nil {
		tx.Rollback()

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := new(domain.User)

		if err = rows.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Surname, &user.Sex, &user.Birthday,
			&user.City, &user.Interests, &user.AccessToken, &user.RefreshToken); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (p *userRepository) CompareError(err error, number uint16) bool {
	me, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}

	return me.Number == number
}

type messengerRepository struct {
	conn *clickhouse.Storage
}

func NewMessengerRepository(conn *clickhouse.Storage) *messengerRepository {
	return &messengerRepository{conn: conn}
}

func (m *messengerRepository) CreateChat(masterID, slaveID string) (string, error) {
	var chatID string

	//SELECT id
	//FROM chat
	//WHERE hasAll(participants, [toUUID('786b020e-a126-4950-abfa-68cd7b6c5d6e'), toUUID('fd883849-4692-44d5-97f6-4d7ca403a207')])

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

func (m *messengerRepository) SendMessages(shardID int, userID, chatID string, messages []*domain.ShortMessage) error {
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
	conns *wstransport.Conns
}

func NewWSPoolRepository(conns *wstransport.Conns) *wsPoolRepository {
	return &wsPoolRepository{
		conns: conns,
	}
}

func (w *wsPoolRepository) AddConnection(userID string, conn net.Conn) {
	w.conns.Add(userID, conn)
}

func (w *wsPoolRepository) RemoveConnection(userID string, conn net.Conn) {
	w.conns.Remove(userID, conn)
}

type cacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) *cacheRepository {
	return &cacheRepository{client: client}
}

func (c *cacheRepository) DoesUserExist(ctx context.Context, userID string) (bool, error) {
	err := c.client.Get(ctx, userID).Err()

	switch err {
	case nil:
		return true, nil
	case redis.Nil:
		return false, nil
	default:
		return false, err
	}
}

func (c *cacheRepository) GetAllUsers(ctx context.Context) ([]string, error) {
	return c.client.Keys(ctx, "*").Result()
}

func (c *cacheRepository) Persist(ctx context.Context, userID string) error {
	return c.client.Set(ctx, userID, true, 0).Err()
}
