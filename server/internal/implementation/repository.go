package implementation

import (
	"context"
	"database/sql"
	"social-network/internal/domain"
)

type profileRepository struct {
	conn *sql.DB
}

func NewProfileRepository(conn *sql.DB) *profileRepository {
	return &profileRepository{conn: conn}
}

func (p *profileRepository) Persist(ctx context.Context, profile *domain.Profile) error {
	stmt, err := p.conn.Prepare(`
		INSERT 
			INTO user (login, password, name, surname, birthday, sex, city, interests) 
		VALUES
		    ( ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	if _, err = stmt.Exec(
		profile.Login,
		profile.Password,
		profile.Name,
		profile.Surname,
		profile.Birthday,
		profile.Sex,
		profile.City,
		profile.Interests); err != nil {
		return err
	}

	return stmt.Close()
}
