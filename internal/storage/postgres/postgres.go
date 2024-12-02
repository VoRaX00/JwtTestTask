package postgres

import (
	models2 "JwtTestTask/internal/domain/models"
	"JwtTestTask/internal/storage"
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

func (s *Storage) CreateToken(token models2.RefreshToken) error {
	query := fmt.Sprintf("INSERT INTO refresh_tokens (user_id, refresh_token_hash, ip, expires_at) VALUES ($1, $2, $3, $4) RETURNING id")
	err := s.DB.QueryRow(query, token.UserId, token.RefreshTokenHash, token.Ip, token.ExpiresAt).Err()
	return err
}

func (s *Storage) GetRefreshToken(refreshTokenHash string) (*models2.RefreshToken, error) {
	const op = "storage.auth.RefreshToken"
	var token models2.RefreshToken
	query := `SELECT ip, user_id, created_at, expires_at FROM refresh_tokens WHERE refresh_token_hash = $1`

	err := s.DB.QueryRow(query, refreshTokenHash).Scan(&token.Ip, &token.UserId, &token.CreatedAt, &token.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.TokenNotFound)
		}
		return nil, err
	}
	return &token, nil
}

func (s *Storage) RefreshToken(newTokenHash, tokenHash, ipClient string, ttl time.Duration) error {
	expiresAt := time.Now().Add(ttl)
	query := `UPDATE refresh_tokens SET refresh_token_hash=$1, expires_at=$2, ip=$3, created_at=CURRENT_TIMESTAMP WHERE refresh_token_hash=$4`
	_, err := s.DB.Exec(query, newTokenHash, expiresAt, ipClient, tokenHash)
	return err
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

func (s *Storage) CreateUser(user models2.User) (string, error) {
	query := fmt.Sprintf("INSERT INTO Users (id, email, password_hash) VALUES ($1, $2, $3) RETURNING id")
	row := s.DB.QueryRow(query, user.Id, user.Email, user.Password)
	if err := row.Scan(&user.Id); err != nil {
		return "", err
	}
	return user.Id, nil
}
