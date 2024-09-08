package repositories

import (
	"JwtTestTask/models"
	"time"
)

type IUserRepository interface {
	Create(user models.User) (string, error)
}

type ITokenRepository interface {
	Create(token models.RefreshToken) error
	RefreshTokens(newTokenHash, tokenHash string, ttl time.Duration) error
	GetUserEmail(token string) (string, error)
}
