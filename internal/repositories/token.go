package repositories

import (
	"JwtTestTask/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
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

func (r *TokenRepository) RefreshTokens(newTokenHash, tokenHash string, ttl time.Duration) error {
	var userId string
	var expiresAt time.Time
	var revoked bool
	query := fmt.Sprintf(`SELETE user_id, expires_at, revoked FROM refresh_tokens WHERE refresh_token_hash=$1`)
	err := r.db.QueryRow(query, tokenHash).Scan(&userId, &expiresAt, &revoked)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("refresh token not found")
		}
		return err
	}

	if revoked {
		return errors.New("token is revoked")
	}
	if time.Now().After(expiresAt) {
		return errors.New("token is expired")
	}

	expiresAt = time.Now().Add(ttl)
	query = fmt.Sprintf("UPDATE refresh_tokens SET refresh_token_hash=$1, expires_at=$2, created_at=CURRENT_TIMESTAMP WHERE user_id=$3")
	_, err = r.db.Exec(query, newTokenHash, expiresAt, userId)
	return err
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
