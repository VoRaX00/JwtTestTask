package repositories

import (
	"JwtTestTask/models"
)

type IUserRepository interface {
	Create(user models.User) (string, error)
}

type ITokenRepository interface {
	Create(token models.RefreshToken) error
	RefreshTokens(tokenHash, ipClient string) (string, error)
	GetUserEmail(token string) (string, error)
}
