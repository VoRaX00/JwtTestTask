package postgres

import (
	"JwtTestTask/internal/storage"
	"JwtTestTask/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type Storage struct {
	DB *sqlx.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"
	db, err := sqlx.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{DB: db}, nil
}

func (s *Storage) CreateToken(token models.RefreshToken) error {
	query := fmt.Sprintf("INSERT INTO refresh_tokens (user_id, refresh_token_hash, ip, expires_at) VALUES ($1, $2, $3, $4) RETURNING id")
	err := s.DB.QueryRow(query, token.UserId, token.RefreshTokenHash, token.Ip, token.ExpiresAt).Err()
	return err
}

func (s *Storage) RefreshToken(refreshTokenHash string) (*models.RefreshToken, error) {
	const op = "storage.auth.RefreshToken"
	var token models.RefreshToken
	query := `SELECT ip, created_at, expires_at FROM refresh_tokens WHERE refresh_token_hash = $1`

	err := s.DB.QueryRow(query, refreshTokenHash).Scan(&token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.TokenNotFound)
		}
		return nil, err
	}
	return &token, nil
}

func (s *Storage) RefreshTokens(newTokenHash, tokenHash, ipClient string, ttl time.Duration) (string, error) {
	var userId string
	var expiresAt time.Time
	var ip string
	query := fmt.Sprintf(`SELECT user_id, expires_at, ip FROM refresh_tokens WHERE refresh_token_hash=$1`)

	err := s.DB.QueryRow(query, tokenHash).Scan(&userId, &expiresAt, &ip)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.TokenNotFound
		}
		return "", err
	}

	if time.Now().After(expiresAt) {
		return "", errors.New("token is expired")
	}

	expiresAt = time.Now().Add(ttl)
	query = fmt.Sprintf("UPDATE refresh_tokens SET refresh_token_hash=$1, expires_at=$2, ip=$3, created_at=CURRENT_TIMESTAMP WHERE user_id=$4")
	_, err = s.DB.Exec(query, newTokenHash, expiresAt, ipClient, userId)
	if ip != ipClient {
		return ipClient, err
	}
	return "", err
}

func (s *Storage) GetUserEmail(token string) (string, error) {
	var userId string
	query := `SELECT user_id FROM refresh_tokens WHERE refresh_token_hash = $1`
	err := s.DB.Get(&userId, query, token)
	if err != nil {
		return "", err
	}

	var email string
	query = `SELECT email FROM users WHERE id = $1`
	err = s.DB.Get(&email, query, userId)
	if err != nil {
		return "", err
	}
	return email, nil
}

func (s *Storage) CreateUser(user models.User) (string, error) {
	query := fmt.Sprintf("INSERT INTO Users (id, email, password_hash) VALUES ($1, $2, $3) RETURNING id")
	row := s.DB.QueryRow(query, user.Id, user.Email, user.Password)
	if row.Scan(&user.Id) != nil {
		return "", nil
	}
	return user.Id, nil
}
