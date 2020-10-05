package domain

import "context"

type ProfileRepository interface {
	Persist(ctx context.Context, profile *Profile) error
	//GetProfileByLogin(ctx context.Context, login string) error
}
