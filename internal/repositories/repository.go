package repositories

import "github.com/jmoiron/sqlx"

type Repository struct {
	IUserRepository
	ITokenRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		IUserRepository:  NewUsersRepository(db),
		ITokenRepository: NewTokenRepository(db),
	}
}
