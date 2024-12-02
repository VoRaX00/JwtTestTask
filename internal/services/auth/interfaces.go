package auth

import (
	models2 "JwtTestTask/internal/domain/models"
	"time"
)

type UserProvider interface {
	CreateUser(user models2.User) (string, error)
}

type TokenProvider interface {
	CreateToken(token models2.RefreshToken) error
	GetRefreshToken(refreshTokenHash string) (*models2.RefreshToken, error)
	RefreshToken(newTokenHash, tokenHash, ipClient string, ttl time.Duration) error
	GetUserEmail(token string) (string, error)
}
