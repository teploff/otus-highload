package implementation

import (
	"auth/internal/domain"
	"context"
	"database/sql"
	"strings"
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

func (p *userRepository) GetByAnthroponym(tx *sql.Tx, anthroponym, userID string, limit, offset int) ([]*domain.User, int, error) {
	var (
		rows  *sql.Rows
		count int
		err   error
	)

	users := make([]*domain.User, 0, 100)
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

		if err = rows.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Surname, &user.Sex, &user.Birthday,
			&user.City, &user.Interests, &user.AccessToken, &user.RefreshToken); err != nil {
			tx.Rollback()

			return nil, 0, err
		}

		users = append(users, user)
	}

	if err = tx.QueryRow(`SELECT FOUND_ROWS()`).Scan(&count); err != nil {
		tx.Rollback()

		return nil, 0, err
	}

	return users, count, nil
}

func (p *userRepository) GetByIDs(tx *sql.Tx, ids []string) ([]*domain.User, error) {
	users := make([]*domain.User, 0, 32)

	sqlStr := "SELECT id, email, password, name, surname, sex, birthday, city, interests, access_token, refresh_token" +
		" FROM user WHERE id IN ("
	vals := make([]interface{}, 0, len(ids))

	for _, id := range ids {
		sqlStr += "?,"
		vals = append(vals, id)
	}

	//trim the last ,
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	// add ) with limit and offset
	sqlStr += ")"

	//prepare the statement
	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		tx.Rollback()

		return nil, err
	}

	rows, err := stmt.Query(vals...)
	if err != nil {
		tx.Rollback()

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := new(domain.User)

		if err = rows.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Surname,
			&user.Sex, &user.Birthday, &user.City, &user.Interests, &user.AccessToken, &user.RefreshToken); err != nil {
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
