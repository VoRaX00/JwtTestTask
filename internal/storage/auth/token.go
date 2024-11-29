package auth

import (
	"JwtTestTask/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	query := fmt.Sprintf("INSERT INTO refresh_tokens (user_id, refresh_token_hash, ip, expires_at) VALUES ($1, $2, $3, $4) RETURNING id")
	err := r.db.QueryRow(query, token.UserId, token.RefreshTokenHash, token.Ip, token.ExpiresAt).Err()
	return err
}

func (r *TokenRepository) RefreshTokens(newTokenHash, tokenHash, ipClient string, ttl time.Duration) (string, error) {
	var userId string
	var expiresAt time.Time
	var ip string
	query := fmt.Sprintf(`SELECT user_id, expires_at, ip FROM refresh_tokens WHERE refresh_token_hash=$1`)

	err := r.db.QueryRow(query, tokenHash).Scan(&userId, &expiresAt, &ip)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("refresh token not found")
		}
		return "", err
	}

	if time.Now().After(expiresAt) {
		return "", errors.New("token is expired")
	}

	expiresAt = time.Now().Add(ttl)
	query = fmt.Sprintf("UPDATE refresh_tokens SET refresh_token_hash=$1, expires_at=$2, ip=$3, created_at=CURRENT_TIMESTAMP WHERE user_id=$4")
	_, err = r.db.Exec(query, newTokenHash, expiresAt, ipClient, userId)
	if ip != ipClient {
		return ipClient, err
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
