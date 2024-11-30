package auth

import (
	"JwtTestTask/internal/services"
	"JwtTestTask/models"
)

type IUserService interface {
	Create(user models.User) (string, error)
	generatePasswordHash(password string) string
}

type ITokenService interface {
	GenerateTokens(userId, ipClient string) (map[string]string, error)
	SendMessageEmail(email, message string) error
	RefreshTokens(token services.Tokens, ipClient string) (map[string]string, error)
	GetUserEmail(token string) (string, error)
	create(token, userId, ipClient string) error
}