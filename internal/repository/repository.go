package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	IAuthRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		IAuthRepository: NewAuthRepository(db),
	}
}
