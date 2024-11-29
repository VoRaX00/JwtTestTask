package storage

import (
	"JwtTestTask/internal/storage/auth"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Token *auth.TokenRepository
	User  *auth.UsersRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Token: auth.NewTokenRepository(db),
		User:  auth.NewUsersRepository(db),
	}
}
