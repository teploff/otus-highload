package implementation

import (
	"context"
	"database/sql"
	uuid "github.com/satori/go.uuid"
	"net"
	"social-network/internal/domain"
	wstransport "social-network/internal/transport/ws"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

const (
	friendshipNonameStatus    = "noname"
	friendshipExpectedStatus  = "expected"
	friendshipConfirmedStatus = "confirmed"
	friendshipAcceptedStatus  = "accepted"
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

func (p *userRepository) GetByAnthroponym(tx *sql.Tx, anthroponym, userID string, limit, offset int) ([]*domain.User, int, error) {
	var (
		rows  *sql.Rows
		count int
		err   error
	)

	friendships := make([]*domain.FriendShip, 0, 100)
	users := make([]*domain.User, 0, 100)

	rows, err = tx.Query(`
		SELECT
			master_user_id, slave_user_id, status
		FROM
			user
		JOIN friendship 
			ON user.id = friendship.master_user_id
		WHERE
			user.id = ?
		UNION
		SELECT
			master_user_id, slave_user_id, status
		FROM
			user
		JOIN friendship 
			ON user.id = friendship.slave_user_id
		WHERE
			user.id = ?`, userID, userID)
	if err != nil {
		tx.Rollback()

		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		friendship := new(domain.FriendShip)

		if err = rows.Scan(&friendship.MasterUserID, &friendship.SlaveUserID, &friendship.Status); err != nil {
			tx.Rollback()

			return nil, 0, err
		}

		friendships = append(friendships, friendship)
	}

	strs := strings.Split(anthroponym, " ")

	if len(strs) > 1 {
		rows, err = tx.Query(`
			SELECT
			    SQL_CALC_FOUND_ROWS id, email, password, name, surname, sex, birthday, city, interests, access_token, refresh_token
			FROM
		    	user
			WHERE 
				(name LIKE ? AND surname LIKE ?) OR (name LIKE ? AND surname LIKE ?) AND id != ?
			ORDER BY surname
			LIMIT ? OFFSET ?`, strs[0]+"%", strs[1]+"%", strs[1]+"%", strs[0]+"%", userID, limit, offset)
	} else {
		rows, err = tx.Query(`
			SELECT
				SQL_CALC_FOUND_ROWS id, email, password, name, surname, sex, birthday, city, interests, access_token, refresh_token
			FROM
		    	user
			WHERE 
				name LIKE ? OR surname LIKE ? AND id != ?
			ORDER BY surname
			LIMIT ? OFFSET ?`, strs[0]+"%", strs[0]+"%", userID, limit, offset)
	}

	if err != nil {
		tx.Rollback()

		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		user := new(domain.User)
		user.FriendshipStatus = friendshipNonameStatus

		if err = rows.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Surname, &user.Sex, &user.Birthday,
			&user.City, &user.Interests, &user.AccessToken, &user.RefreshToken); err != nil {
			tx.Rollback()

			return nil, 0, err
		}

		for _, friendship := range friendships {
			if friendship.MasterUserID == user.ID {
				switch friendship.Status {
				case friendshipExpectedStatus:
					user.FriendshipStatus = friendshipConfirmedStatus
				default:
					user.FriendshipStatus = friendshipAcceptedStatus
				}
			} else if friendship.SlaveUserID == user.ID {
				user.FriendshipStatus = friendship.Status
			}
		}

		users = append(users, user)
	}

	if err = tx.QueryRow(`SELECT FOUND_ROWS()`).Scan(&count); err != nil {
		tx.Rollback()

		return nil, 0, err
	}

	return users, count, nil
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

func (m *messengerRepository) GetChatWithCompanion(tx *sql.Tx, masterID, slaveID string) (*domain.Chat, error) {
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

func (m *messengerRepository) GetChatAsParticipant(tx *sql.Tx, userID string) (*domain.Chat, error) {
	var chat domain.Chat

	err := tx.QueryRow(`
		SELECT
			chat.id, chat.create_time
		FROM user
		JOIN user_chat
			ON user.id = user_chat.user_id
		JOIN chat
			ON user_chat.chat_id = chat.id
		WHERE user.id = ?`, userID).Scan(&chat.ID, &chat.CreateTime)
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

func (m *messengerRepository) GetCountMessages(tx *sql.Tx, chatID string) (int, error) {
	var count int
	//select count(*) from (select max(create_time) from message where chat_id = "188a4d72-64a7-4dd4-beae-df8ca11fce70" group by id) AS a;
	err := tx.QueryRow(`
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
		tx.Rollback()

		return 0, err
	}

	return count, nil
}

func (m *messengerRepository) GetMessages(tx *sql.Tx, chatID string, limit, offset int) ([]*domain.Message, error) {
	messages := make([]*domain.Message, 0, 10)

	rows, err := tx.Query(`
		SELECT
			id, text, status, user_id, max(create_time) as create_time
		FROM message
		WHERE chat_id = ?
		GROUP BY id, text, status, user_id
		LIMIT ? OFFSET ?`, chatID, limit, offset)
	if err != nil {
		tx.Rollback()

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message domain.Message

		if err = rows.Scan(&message.ID, &message.Text, &message.Status, &message.UserID, &message.CreateTime); err != nil {
			tx.Rollback()

			return nil, err
		}

		messages = append(messages, &message)
	}

	return messages, nil
}

type socialRepository struct {
	conn *sql.DB
}

func NewSocialRepository(conn *sql.DB) *socialRepository {
	return &socialRepository{conn: conn}
}

func (s *socialRepository) GetTx(ctx context.Context) (*sql.Tx, error) {
	return s.conn.BeginTx(ctx, nil)
}

func (s *socialRepository) CommitTx(tx *sql.Tx) error {
	return tx.Commit()
}

func (s *socialRepository) CreateFriendship(tx *sql.Tx, masterUserID string, slaveUsersID []string) error {
	sqlStr := "INSERT INTO friendship (master_user_id, slave_user_id, status, create_time) VALUES "
	vals := make([]interface{}, 0, len(slaveUsersID)*4) // 4 - count cells: master_user_id, slave_user_id ...

	for _, id := range slaveUsersID {
		sqlStr += "( ?, ?, ?, ?),"
		vals = append(vals, masterUserID, id, "expected", time.Now().UTC())
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

func (s *socialRepository) ConfirmFriendship(tx *sql.Tx, userID string, friendsID []string) error {
	for _, friendID := range friendsID {
		_, err := tx.Exec(`
		UPDATE 
			friendship
		SET
		    status = ?
		WHERE
		    master_user_id = ? AND slave_user_id = ?`, friendshipAcceptedStatus, friendID, userID)
		if err != nil {
			tx.Rollback()

			return err
		}
	}

	return nil
}

func (s *socialRepository) RejectFriendship(tx *sql.Tx, userID string, friendsID []string) error {
	sqlStr := "DELETE FROM friendship WHERE (master_user_id,slave_user_id) IN ("
	vals := make([]interface{}, 0, len(friendsID)*2) // 4 - count cells: master_user_id, slave_user_id ...

	for _, friendID := range friendsID {
		sqlStr += "( ?, ?),"
		vals = append(vals, friendID, userID)
	}

	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	// add )
	sqlStr += ")"

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

func (s *socialRepository) BreakFriendship(tx *sql.Tx, userID string, friendsID []string) error {
	for _, friendID := range friendsID {
		_, err := tx.Exec(`
		DELETE
			FROM friendship
		WHERE
		    (master_user_id = ? AND slave_user_id = ?) OR (master_user_id = ? AND slave_user_id = ?) AND status = ?`,
			friendID, userID, userID, friendID, friendshipAcceptedStatus)
		if err != nil {
			tx.Rollback()

			return err
		}
	}

	return nil
}

func (s *socialRepository) GetFriends(tx *sql.Tx, userID string) ([]*domain.User, error) {
	var (
		rows *sql.Rows
		err  error
	)

	users := make([]*domain.User, 0, 100)

	rows, err = tx.Query(`
		SELECT
			user.id, user.email, user.password, user.name, user.surname, user.sex, user.birthday, user.city, user.interests, user.access_token, user.refresh_token
		FROM
			user
		JOIN friendship 
			ON user.id = friendship.master_user_id
		WHERE
			friendship.slave_user_id = ? and friendship.status = ?
		UNION
		SELECT
			user.id, user.email, user.password, user.name, user.surname, user.sex, user.birthday, user.city, user.interests, user.access_token, user.refresh_token
		FROM
			user
		JOIN friendship 
			ON user.id = friendship.slave_user_id
		WHERE
			friendship.master_user_id = ? and friendship.status = ?`, userID, friendshipAcceptedStatus, userID, friendshipAcceptedStatus)
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

func (s *socialRepository) GetFollowers(tx *sql.Tx, userID string) ([]*domain.User, error) {
	var (
		rows *sql.Rows
		err  error
	)

	users := make([]*domain.User, 0, 100)

	rows, err = tx.Query(`
		SELECT
			user.id, user.email, user.password, user.name, user.surname, user.sex, user.birthday, user.city, user.interests, user.access_token, user.refresh_token
		FROM
			user
		JOIN friendship
			ON user.id = friendship.master_user_id
		WHERE
			friendship.slave_user_id = ? and friendship.status = ?`, userID, friendshipExpectedStatus)
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
