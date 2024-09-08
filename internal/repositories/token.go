package repositories

import (
	"JwtTestTask/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type TokenRepository struct {
	db *sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) *TokenRepository {
	return &TokenRepository{
		db: db,
	}
}

func (r *TokenRepository) Create(token models.RefreshToken) error {
	query := fmt.Sprintf("INSERT INTO refresh_tokens (userId, refresh_token_hash, expires_at) VALUES ($1, $2, $3) RETURNING id")
	err := r.db.QueryRow(query, token.UserId, token.RefreshTokenHash, token.ExpiresAt).Err()
	return err
}

func (r *TokenRepository) RefreshTokens(tokenHash, ipClient string) (string, error) {
	var token models.RefreshToken
	query := fmt.Sprintf("SELETE * FROM refresh_tokens WHERE refresh_token_hash=$1")
	err := r.db.Get(&token, query, tokenHash)
	if err != nil {
		return "", err
	}
	return "", err
}

func (r *TokenRepository) GetUserEmail(token string) (string, error) {
	var userId string
	query := `SELECT user_id FROM refresh_tokens WHERE refresh_token_hash = $1`
	err := r.db.Get(&userId, query, token)
	if err != nil {
		return "", err
	}

	var email string
	query = `SELECT email FROM users WHERE id = $1`
	err = r.db.Get(&email, query, userId)
	if err != nil {
		return "", err
	}
	return email, nil
}
