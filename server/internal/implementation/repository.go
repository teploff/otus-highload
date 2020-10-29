package implementation

import (
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"social-network/internal/domain"
	"time"
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
