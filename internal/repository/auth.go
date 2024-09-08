package repository

import (
	"JwtTestTask/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) Get(user models.User) (string, error) {
	var id string
	query := fmt.Sprintf("SELECT id FROM Users WHERE email=$1 AND password_hash=$2")
	err := r.db.Get(&id, query, user.Email, user.Password)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *AuthRepository) Create(user models.User) (string, error) {
	query := fmt.Sprintf("INSERT INTO Users (id, email, password_hash) VALUES ($1, $2, $3) RETURNING id")
	row := r.db.QueryRow(query, user.Id, user.Email, user.Password)
	if row.Scan(&user.Id) != nil {
		return "", nil
	}
	return user.Id, nil
}

func (r *AuthRepository) AddToken(token models.RefreshToken) error {
	expiresAt := time.Now().Add(token.TTL)
	query := fmt.Sprintf("INSERT INTO refresh_tokens (userId, refresh_token_hash, expires_at) VALUES ($1, $2, $3) RETURNING id")
	err := r.db.QueryRow(query, token.UserId, token.TokenHash, expiresAt).Err()
	return err
}
