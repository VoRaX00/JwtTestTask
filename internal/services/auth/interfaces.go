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
	RefreshToken(refreshTokenHash string) (*models.RefreshToken, error)
	RefreshTokens(newTokenHash, tokenHash, ipClient string, ttl time.Duration) (string, error)
	GetUserEmail(token string) (string, error)
}
