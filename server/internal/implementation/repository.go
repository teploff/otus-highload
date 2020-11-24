package implementation

import (
	"context"
	"database/sql"
	uuid "github.com/satori/go.uuid"
	"net"
	"social-network/internal/domain"
	wstransport "social-network/internal/transport/ws"
	"time"

	"github.com/go-sql-driver/mysql"
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
			tx.Rollback()

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
			tx.Rollback()

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

func (m *messengerRepository) CreateChat(tx *sql.Tx, masterID, slaveID string) (string, error) {
	var chatID string
	err := tx.QueryRow(`
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
		tx.Rollback()

		return "", err
	}

	_, err = tx.Exec(`
		INSERT 
			INTO chat (id, create_time) 
		VALUES
			(?, ?)`, chatID, time.Now().UTC())
	if err != nil {
		tx.Rollback()

		return "", err
	}

	_, err = tx.Exec(`
		INSERT 
			INTO user_chat (user_id, chat_id) 
		VALUES
			(?, ?)`, masterID, chatID)
	if err != nil {
		tx.Rollback()

		return "", err
	}

	_, err = tx.Exec(`
		INSERT 
			INTO user_chat (user_id, chat_id) 
		VALUES
			(?, ?)`, slaveID, chatID)
	if err != nil {
		tx.Rollback()

		return "", err
	}

	return chatID, nil
}

func (m *messengerRepository) GetCountChats(tx *sql.Tx, userID string) (int, error) {
	var count int

	err := tx.QueryRow(`
		SELECT 
			count(*)
		FROM user_chat
		JOIN chat
		    ON user_chat.chat_id = chat.id
		WHERE user_chat.user_id = ?`, userID).Scan(&count)
	if err != nil {
		tx.Rollback()

		return 0, err
	}

	return count, nil
}

func (m *messengerRepository) GetChat(tx *sql.Tx, masterID, slaveID string) (*domain.Chat, error) {
	var chat domain.Chat

	err := tx.QueryRow(`
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
		tx.Rollback()

		return nil, err
	}

	return &chat, nil
}

func (m *messengerRepository) GetChats(tx *sql.Tx, userID string, limit, offset int) ([]*domain.Chat, error) {
	chats := make([]*domain.Chat, 0, 10)

	rows, err := tx.Query(`
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
		tx.Rollback()

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
			tx.Rollback()

			return nil, err
		}

		chatRows = append(chatRows, &row)
	}
	rows.Close()

	for _, chat := range chatRows {
		rows, err = tx.Query(`
		SELECT
			user.id, user.name, user.surname
		FROM user_chat
		JOIN user
			ON user_chat.user_id = user.id
		WHERE user_chat.chat_id = ? AND user_chat.user_id != ?`, chat.id, userID)
		if err != nil {
			tx.Rollback()

			return nil, err
		}

		for rows.Next() {
			var user struct {
				ID      string
				Name    string
				Surname string
			}
			if err = rows.Scan(&user.ID, &user.Name, &user.Surname); err != nil {
				tx.Rollback()

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

func (m *messengerRepository) SendMessages(tx *sql.Tx, userID, chatID string, messages []*domain.ShortMessage) error {
	sqlStr := "INSERT INTO message (text, status, create_time, user_id, chat_id) VALUES "
	vals := make([]interface{}, 0, len(messages)*6)

	for _, msg := range messages {
		sqlStr += "( ?, ?, ?, ?, ?),"
		vals = append(vals, msg.Text, msg.Status, time.Now().UTC(), userID, chatID)
	}

	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]

	//prepare the statement
	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		tx.Rollback()

		return err
	}

	//format all vals at once
	if _, err = stmt.Exec(vals...); err != nil {
		tx.Rollback()

		return err
	}

	return nil
}

func (m *messengerRepository) GetCountMessages(tx *sql.Tx, userID, chatID string) (int, error) {
	var count int

	err := tx.QueryRow(`
		SELECT 
			count(*)
		FROM user_chat
		JOIN chat
		    ON user_chat.chat_id = chat.id
		JOIN message
			ON user_chat.chat_id = message.chat_id
		WHERE user_chat.user_id = ? AND user_chat.chat_id = ?`, userID, chatID).Scan(&count)
	if err != nil {
		tx.Rollback()

		return 0, err
	}

	return count, nil
}

func (m *messengerRepository) GetMessages(tx *sql.Tx, userID, chatID string, limit, offset int) ([]*domain.Message, error) {
	messages := make([]*domain.Message, 0, 10)

	rows, err := tx.Query(`
		SELECT
			id, text, max(status), create_time
		FROM message
		WHERE user_id = ? AND chat_id = ?
		GROUP BY id
		LIMIT ? OFFSET ?`, userID, chatID, limit, offset)
	if err != nil {
		tx.Rollback()

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message domain.Message

		if err = rows.Scan(&message.ID, &message.Text, &message.Status, &message.CreateTime); err != nil {
			tx.Rollback()

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
