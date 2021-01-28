package implementation

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net"
	"social/internal/domain"
	"social/internal/infrastructure/cache"
	"sync"
	"time"
)

const (
	friendshipNonameStatus    = "noname"
	friendshipExpectedStatus  = "expected"
	friendshipConfirmedStatus = "confirmed"
	friendshipAcceptedStatus  = "accepted"
)

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

func (s *socialRepository) GetFriends(tx *sql.Tx, userID string) ([]string, error) {
	usersID := make([]string, 0, 100)

	rows, err := tx.Query(`
		SELECT
			slave_user_id as id
		FROM
			friendship
		WHERE
			master_user_id = ? and status = ?
		UNION
		SELECT
			slave_user_id as id
		FROM
			friendship
		WHERE
			slave_user_id = ? and status = ?`, userID, friendshipAcceptedStatus, userID, friendshipAcceptedStatus)
	if err != nil {
		tx.Rollback()

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string

		if err = rows.Scan(&id); err != nil {
			tx.Rollback()

			return nil, err
		}

		usersID = append(usersID, id)
	}

	return usersID, nil
}

func (s *socialRepository) GetFollowers(tx *sql.Tx, userID string) ([]string, error) {
	usersID := make([]string, 0, 100)

	rows, err := tx.Query(`
		SELECT
			master_user_id as id
		FROM
			friendship
		WHERE
		    slave_user_id = ? and status = ?`, userID, friendshipExpectedStatus)
	if err != nil {
		tx.Rollback()

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID string

		if err = rows.Scan(&userID); err != nil {
			tx.Rollback()

			return nil, err
		}

		usersID = append(usersID, userID)
	}

	return usersID, nil
}

func (s *socialRepository) GetUserFriendships(tx *sql.Tx, userID string) ([]*domain.FriendShip, error) {
	friendships := make([]*domain.FriendShip, 0, 100)

	rows, err := tx.Query(`
		SELECT
			master_user_id, slave_user_id, status
		FROM
			friendship
		WHERE
			master_user_id = ? OR slave_user_id = ?`, userID, userID)
	if err != nil {
		tx.Rollback()

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		friendship := new(domain.FriendShip)

		if err = rows.Scan(&friendship.MasterUserID, &friendship.SlaveUserID, &friendship.Status); err != nil {
			tx.Rollback()

			return nil, err
		}

		friendships = append(friendships, friendship)
	}

	return friendships, nil
}

func (s *socialRepository) GetNews(tx *sql.Tx, friends []*domain.User, limit, offset int) ([]*domain.News, int, error) {
	friendsID := make([]string, 0, len(friends))

	for _, id := range friendsID {
		friendsID = append(friendsID, id)
	}

	var count int
	news := make([]*domain.News, 0, 100)

	sqlStr := "SELECT SQL_CALC_FOUND_ROWS id, owner_id, content, create_time FROM news WHERE owner_id IN ("
	vals := make([]interface{}, 0, len(friendsID))

	for _, id := range friendsID {
		sqlStr += "?,"
		vals = append(vals, id)
	}
	vals = append(vals, limit, offset)

	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	// add ) with limit and offset
	sqlStr += ") ORDER BY news.create_time DESC LIMIT ? OFFSET ?"

	//prepare the statement
	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		tx.Rollback()

		return nil, 0, err
	}

	rows, err := stmt.Query(vals...)
	if err != nil {
		tx.Rollback()

		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var ownerID string
		n := new(domain.News)

		if err = rows.Scan(&n.ID, &ownerID, &n.Content, &n.CreateTime); err != nil {
			tx.Rollback()

			return nil, 0, err
		}

		for _, u := range friends {
			if u.ID == ownerID {
				n.Owner.Name = u.Name
				n.Owner.Surname = u.Surname
				n.Owner.Sex = u.Sex

				break
			}
		}

		news = append(news, n)
	}

	if err = tx.QueryRow(`SELECT FOUND_ROWS()`).Scan(&count); err != nil {
		tx.Rollback()

		return nil, 0, err
	}

	return news, count, nil
}

func (s *socialRepository) PublishNews(tx *sql.Tx, userID string, news []*domain.News) error {
	sqlStr := "INSERT INTO news (id, owner_id, content, create_time) VALUES "
	vals := make([]interface{}, 0, len(news)*4) // 4 - count cells: id, owner_id, content, create_time

	for _, n := range news {
		sqlStr += "( ?, ?, ?, ?),"
		vals = append(vals, n.ID, userID, n.Content, n.CreateTime)
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

type socialCacheRepository struct {
	pool           *cache.Pool
	friendsDBIndex int
	NewsDBIndex    int
}

func NewCacheRepository(pool *cache.Pool) *socialCacheRepository {
	return &socialCacheRepository{
		pool:           pool,
		friendsDBIndex: 1,
		NewsDBIndex:    2,
	}
}

func (s *socialCacheRepository) PersistFriend(ctx context.Context, userID string, friendsID []string) error {
	conn, err := s.pool.GetConnByIndexDB(s.friendsDBIndex)
	if err != nil {
		return err
	}

	var friends []string
	result, err := conn.Get(ctx, userID).Result()
	switch err {
	case nil:
		if err = json.Unmarshal([]byte(result), &friends); err != nil {
			return fmt.Errorf("cannot unmarshal friends id, %w", err)
		}
	case redis.Nil:
		friends = make([]string, 0, 1)
	default:
		return err
	}
	friends = append(friends, friendsID...)

	data, err := json.Marshal(friends)
	if err != nil {
		return fmt.Errorf("cannot marshal friends id, %w", err)
	}

	return conn.Set(ctx, userID, data, 0).Err()
}

func (s *socialCacheRepository) DeleteFriend(ctx context.Context, userID string, friendsID []string) error {
	conn, err := s.pool.GetConnByIndexDB(s.friendsDBIndex)
	if err != nil {
		return err
	}

	var friends []string
	result, err := conn.Get(ctx, userID).Result()
	switch err {
	case nil:
		if err = json.Unmarshal([]byte(result), &friends); err != nil {
			return fmt.Errorf("cannot unmarshal friends id, %w", err)
		}
	case redis.Nil:
		return fmt.Errorf("friends are absent")
	default:
		return err
	}

	for i := 0; i < len(friends); i++ {
		friend := friends[i]
		for _, remFriend := range friendsID {
			if friend == remFriend {
				friends = append(friends[:i], friends[i+1:]...)
				i-- // Important: decrease index
				break
			}
		}
	}

	data, err := json.Marshal(friends)
	if err != nil {
		return fmt.Errorf("cannot marshal friends id, %w", err)
	}

	return conn.Set(ctx, userID, data, 0).Err()
}

func (s *socialCacheRepository) RetrieveFriendsID(ctx context.Context, userID string) ([]string, error) {
	conn, err := s.pool.GetConnByIndexDB(s.friendsDBIndex)
	if err != nil {
		return nil, err
	}

	var friendsID []string
	result, err := conn.Get(ctx, userID).Result()
	switch err {
	case nil:
		if err = json.Unmarshal([]byte(result), &friendsID); err != nil {
			return nil, fmt.Errorf("cannot unmarshal friends id, %w", err)
		}
	case redis.Nil:
		friendsID = make([]string, 0, 1)
	default:
		return nil, err
	}

	return friendsID, nil
}

func (s *socialCacheRepository) PersistNews(ctx context.Context, userID string, news []*domain.News) error {
	conn, err := s.pool.GetConnByIndexDB(s.NewsDBIndex)
	if err != nil {
		return err
	}

	var n []*domain.News
	result, err := conn.Get(ctx, userID).Result()
	switch err {
	case nil:
		if err = json.Unmarshal([]byte(result), &n); err != nil {
			return fmt.Errorf("cannot unmarshal news, %w", err)
		}
	case redis.Nil:
		n = make([]*domain.News, 0, 1)
	default:
		return err
	}
	n = append(n, news...)

	data, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("cannot marshal news, %w", err)
	}

	return conn.Set(ctx, userID, data, 0).Err()
}

func (s *socialCacheRepository) RetrieveNews(ctx context.Context, userID string) ([]*domain.News, error) {
	conn, err := s.pool.GetConnByIndexDB(s.NewsDBIndex)
	if err != nil {
		return nil, err
	}

	var news []*domain.News
	result, err := conn.Get(ctx, userID).Result()
	switch err {
	case nil:
		if err = json.Unmarshal([]byte(result), &news); err != nil {
			return nil, fmt.Errorf("cannot news, %w", err)
		}
	case redis.Nil:
		news = make([]*domain.News, 0, 1)
	default:
		return nil, err
	}

	return news, nil
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
