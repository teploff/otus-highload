package implementation

import (
	"auth/internal/domain"
	"context"
	"database/sql"
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

func (p *userRepository) CompareError(err error, number uint16) bool {
	me, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}

	return me.Number == number
}
