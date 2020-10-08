package implementation

import (
	"context"
	"database/sql"
	"log"
	"social-network/internal/domain"
	"time"
)

type userRepository struct {
	conn *sql.DB
}

func NewUserRepository(conn *sql.DB) *userRepository {
	return &userRepository{conn: conn}
}

func (p *userRepository) Persist(ctx context.Context, user *domain.User) error {
	stmt, err := p.conn.Prepare(`
		INSERT 
			INTO user (login, password, name, surname, birthday, sex, city, interests) 
		VALUES
		    ( ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Login, user.Password, user.Name, user.Surname, user.Birthday, user.Sex, user.City,
		user.Interests); err != nil {
		return err
	}

	return nil
}

func (p *userRepository) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User

	stmt, err := p.conn.Prepare(`
	SELECT
		id, login, password, name, surname, sex, birthday, city, interests, access_token, refresh_token
	FROM
	     user 
	WHERE 
	      login = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if err = stmt.QueryRowContext(ctx, login).Scan(&user.ID, &user.Login, &user.Password, &user.Name, &user.Surname,
		&user.Sex, &user.Birthday, &user.City, &user.Interests, &user.AccessToken, &user.RefreshToken); err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *userRepository) GetByIDAndRefreshToken(ctx context.Context, id, token string) (*domain.User, error) {
	var user domain.User

	stmt, err := p.conn.Prepare(`
	SELECT
		id, login, password, name, surname, sex, birthday, city, interests, access_token, refresh_token
	FROM
	     user
	WHERE
	      id = ? AND refresh_token = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	log.Println(id, token)
	if err = stmt.QueryRowContext(ctx, id, token).Scan(&user.ID, &user.Login, &user.Password, &user.Name, &user.Surname,
		&user.Sex, &user.Birthday, &user.City, &user.Interests, &user.AccessToken, &user.RefreshToken); err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *userRepository) UpdateByID(ctx context.Context, user *domain.User) error {
	stmt, err := p.conn.Prepare(`
		UPDATE 
			user
		SET
		    login = ?, password = ?, name = ?, surname = ?, birthday = ?, sex = ?, city = ?, interests = ?,
		    access_token = ?, refresh_token = ?, update_time = ?
		WHERE
		    id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Login, user.Password, user.Name, user.Surname, user.Birthday, user.Sex, user.City,
		user.Interests, user.AccessToken, user.RefreshToken, time.Now().UTC(), user.ID); err != nil {
		return err
	}

	return nil
}
