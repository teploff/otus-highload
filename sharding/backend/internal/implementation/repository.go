package implementation

import (
	"backend/internal/domain"
	wstransport "backend/internal/transport/ws"
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net"
	"time"

	"github.com/go-sql-driver/mysql"
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

func (p *userRepository) Persist(user *domain.User) error {
	_, err := p.conn.Exec(`
		INSERT 
			INTO user (email, password, name, surname, birthday, sex, city, interests) 
		VALUES
			( ?, ?, ?, ?, ?, ?, ?, ?)`, user.Email, user.Password, user.Name, user.Surname, user.Birthday, user.Sex,
		user.City, user.Interests)
	if err != nil {
		return err
	}

	return nil
}

func (p *userRepository) GetByID(id string) (*domain.User, error) {
	var user domain.User

	err := p.conn.QueryRow(`
		SELECT
			id, email, password, name, surname, sex, birthday, city, interests, access_token, refresh_token
		FROM
			 user 
		WHERE 
			  id = ?`, id).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Surname,
		&user.Sex, &user.Birthday, &user.City, &user.Interests, &user.AccessToken, &user.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User

	err := p.conn.QueryRow(`
		SELECT
			id, email, password, name, surname, sex, birthday, city, interests, access_token, refresh_token
		FROM
			 user 
		WHERE 
			  email = ?`, email).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Surname,
		&user.Sex, &user.Birthday, &user.City, &user.Interests, &user.AccessToken, &user.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *userRepository) GetByIDAndRefreshToken(id, token string) (*domain.User, error) {
	var user domain.User

	err := p.conn.QueryRow(`
		SELECT
			id, email, password, name, surname, sex, birthday, city, interests, access_token, refresh_token
		FROM
			user
		WHERE
			id = ? AND refresh_token = ?`, id, token).Scan(&user.ID, &user.Email, &user.Password, &user.Name,
		&user.Surname, &user.Sex, &user.Birthday, &user.City, &user.Interests, &user.AccessToken,
		&user.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *userRepository) GetByIDAndAccessToken(id, token string) (*domain.User, error) {
	var user domain.User

	err := p.conn.QueryRow(`
		SELECT
			id, email, password, name, surname, sex, birthday, city, interests, access_token, refresh_token
		FROM
			user
		WHERE
			id = ? AND access_token = ?`, id, token).Scan(&user.ID, &user.Email, &user.Password, &user.Name,
		&user.Surname, &user.Sex, &user.Birthday, &user.City, &user.Interests, &user.AccessToken,
		&user.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *userRepository) UpdateByID(user *domain.User) error {
	_, err := p.conn.Exec(`
		UPDATE 
			user
		SET
		    email = ?, password = ?, name = ?, surname = ?, birthday = ?, sex = ?, city = ?, interests = ?,
		    access_token = ?, refresh_token = ?, update_time = ?
		WHERE
		    id = ?`, user.Email, user.Password, user.Name, user.Surname, user.Birthday, user.Sex,
		user.City, user.Interests, user.AccessToken, user.RefreshToken, time.Now().UTC(), user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *userRepository) GetCount() (int, error) {
	var count int

	if err := p.conn.QueryRow(`SELECT count(*) FROM user`).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (p *userRepository) GetByLimitAndOffsetExceptUserID(userID string, limit, offset int) ([]*domain.User, error) {
	users := make([]*domain.User, 0, 10)

	rows, err := p.conn.Query(`
		SELECT
			id, email, password, name, surname, sex, birthday, city, interests, access_token, refresh_token
		FROM
		    user
		WHERE 
			  id != ?
		ORDER BY create_time
		LIMIT ? OFFSET ?`, userID, limit, offset)
	if err != nil {
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

func (p *userRepository) GetByPrefixOfNameAndSurname(prefix string) ([]*domain.User, error) {
	users := make([]*domain.User, 0, 100)

	rows, err := p.conn.Query(`
		SELECT
			id, email, password, name, surname, sex, birthday, city, interests, access_token, refresh_token
		FROM
		    user
		WHERE 
			  name LIKE ? AND surname LIKE ?
		ORDER BY id`, prefix+"%", prefix+"%")
	if err != nil {
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
	conn *sql.DB
}

func NewMessengerRepository(conn *sql.DB) *messengerRepository {
	return &messengerRepository{conn: conn}
}

func (m *messengerRepository) GetTx(ctx context.Context) (*sql.Tx, error) {
	return m.conn.BeginTx(ctx, nil)
}

func (m *messengerRepository) CommitTx(tx *sql.Tx) error {
	return tx.Commit()
}

func (m *messengerRepository) CreateChat(masterID, slaveID string) (string, error) {
	var chatID string
	err := m.conn.QueryRow(`
		SELECT
			UC1.chat_id
		FROM 
		(
		    SELECT
		    	user_id, chat_id
		    FROM user_chat
		    WHERE user_id = ?
		) UC1
		JOIN (
		    SELECT
		    	user_id, chat_id
		    FROM user_chat
		    WHERE user_id = ?
		) UC2
		    ON UC1.chat_id = UC2.chat_id`, masterID, slaveID).Scan(&chatID)
	switch err {
	case nil:
		return chatID, nil
	case sql.ErrNoRows:
		chatID = uuid.NewV4().String()
	default:
		return "", err
	}

	_, err = m.conn.Exec(`
		INSERT 
			INTO chat (id, create_time) 
		VALUES
			(?, ?)`, chatID, time.Now().UTC())
	if err != nil {
		return "", err
	}

	_, err = m.conn.Exec(`
		INSERT 
			INTO user_chat (user_id, chat_id) 
		VALUES
			(?, ?)`, masterID, chatID)
	if err != nil {
		return "", err
	}

	_, err = m.conn.Exec(`
		INSERT 
			INTO user_chat (user_id, chat_id) 
		VALUES
			(?, ?)`, slaveID, chatID)
	if err != nil {
		return "", err
	}

	return chatID, nil
}

func (m *messengerRepository) GetCountChats(userID string) (int, error) {
	var count int

	err := m.conn.QueryRow(`
		SELECT 
			count(*)
		FROM user_chat
		JOIN chat
		    ON user_chat.chat_id = chat.id
		WHERE user_chat.user_id = ?`, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *messengerRepository) GetChatWithCompanion(masterID, slaveID string) (*domain.Chat, error) {
	var chat domain.Chat

	err := m.conn.QueryRow(`
		SELECT 
			C1.id, C1.create_time
		FROM (
			SELECT
				chat.id, chat.create_time
			FROM user
			JOIN user_chat
				ON user.id = user_chat.user_id
			JOIN chat
				ON user_chat.chat_id = chat.id
			WHERE user.id = ?
		) as C1
		JOIN user_chat AS UC1
			ON C1.id = UC1.chat_id
		where UC1.user_id = ?`, masterID, slaveID).Scan(&chat.ID, &chat.CreateTime)
	if err != nil {
		return nil, err
	}

	return &chat, nil
}

func (m *messengerRepository) GetChatAsParticipant(userID string) (*domain.Chat, error) {
	var chat domain.Chat

	err := m.conn.QueryRow(`
		SELECT
			chat.id, chat.create_time
		FROM user
		JOIN user_chat
			ON user.id = user_chat.user_id
		JOIN chat
			ON user_chat.chat_id = chat.id
		WHERE user.id = ?`, userID).Scan(&chat.ID, &chat.CreateTime)
	if err != nil {
		return nil, err
	}

	return &chat, nil
}

func (m *messengerRepository) GetChats(userID string, limit, offset int) ([]*domain.Chat, error) {
	chats := make([]*domain.Chat, 0, 10)

	rows, err := m.conn.Query(`
		SELECT
			chat.id, chat.create_time
		FROM user
		JOIN user_chat
			ON user.id = user_chat.user_id
		JOIN chat
			ON user_chat.chat_id = chat.id
		WHERE user.id = ?
		ORDER BY chat.id
		LIMIT ? OFFSET ?`, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	type chatRow struct {
		id         string
		createTime time.Time
	}
	chatRows := make([]*chatRow, 0)

	for rows.Next() {
		var row chatRow
		if err = rows.Scan(&row.id, &row.createTime); err != nil {
			return nil, err
		}

		chatRows = append(chatRows, &row)
	}
	rows.Close()

	for _, chat := range chatRows {
		rows, err = m.conn.Query(`
		SELECT
			user.id, user.name, user.surname
		FROM user_chat
		JOIN user
			ON user_chat.user_id = user.id
		WHERE user_chat.chat_id = ? AND user_chat.user_id != ?`, chat.id, userID)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var user struct {
				ID      string
				Name    string
				Surname string
			}
			if err = rows.Scan(&user.ID, &user.Name, &user.Surname); err != nil {
				return nil, err
			}

			exist := false
			for _, c := range chats {
				if c.ID == chat.id {
					c.Participants = append(c.Participants, &domain.Participant{
						ID:      user.ID,
						Name:    user.Name,
						Surname: user.Surname,
					})

					exist = true
				}
			}

			if !exist {
				c := &domain.Chat{
					ID:           chat.id,
					CreateTime:   chat.createTime,
					Participants: make([]*domain.Participant, 0),
				}
				c.Participants = append(c.Participants, &domain.Participant{
					ID:      user.ID,
					Name:    user.Name,
					Surname: user.Surname,
				})

				chats = append(chats, c)
			}
		}
		rows.Close()
	}

	return chats, nil
}

func (m *messengerRepository) SendMessages(shardID int, userID, chatID string, messages []*domain.ShortMessage) error {
	sqlStr := "INSERT INTO message (shard_id, text, status, create_time, user_id, chat_id) VALUES "
	vals := make([]interface{}, 0, len(messages)*6)

	for _, msg := range messages {
		sqlStr += fmt.Sprintf("(%d, ?, ?, ?, ?, ?),", shardID)
		vals = append(vals, msg.Text, msg.Status, time.Now().UTC(), userID, chatID)
	}

	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]

	//prepare the statement
	stmt, err := m.conn.Prepare(sqlStr)
	if err != nil {
		return err
	}

	//format all vals at once
	if _, err = stmt.Exec(vals...); err != nil {
		return err
	}

	return nil
}

func (m *messengerRepository) GetCountMessages(chatID string) (int, error) {
	var count int
	//select count(*) from (select max(create_time) from message where chat_id = "188a4d72-64a7-4dd4-beae-df8ca11fce70" group by id) AS a;
	err := m.conn.QueryRow(`
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

	rows, err := m.conn.Query(`
		SELECT
			id, text, status, user_id, max(create_time) as create_time
		FROM message
		WHERE chat_id = ?
		GROUP BY id, text, status, user_id
		LIMIT ? OFFSET ?`, chatID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message domain.Message

		if err = rows.Scan(&message.ID, &message.Text, &message.Status, &message.UserID, &message.CreateTime); err != nil {
			return nil, err
		}

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

func (c *cacheRepository) GetLadyGagaUsers(ctx context.Context) ([]string, error) {
	return c.client.Keys(ctx, "*").Result()
}

func (c *cacheRepository) Persist(ctx context.Context, userID string) error {
	return c.client.Set(ctx, userID, true, 0).Err()
}
