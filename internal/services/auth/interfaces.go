package auth

import (
	"JwtTestTask/models"
	"time"
)

type UserProvider interface {
	CreateUser(user models.User) (string, error)
}

type TokenProvider interface {
	CreateToken(token models.RefreshToken) error
	GetRefreshToken(refreshTokenHash string) (*models.RefreshToken, error)
	RefreshToken(newTokenHash, tokenHash, ipClient string, ttl time.Duration) error
	GetUserEmail(token string) (string, error)
}
