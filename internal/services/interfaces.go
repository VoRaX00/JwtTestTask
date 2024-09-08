package services

import "JwtTestTask/models"

type IUserService interface {
	Create(user models.User) (string, error)
	SendMessageEmail(email, message string) error
	generatePasswordHash(password string) string
}

type ITokenService interface {
	GenerateTokens(userId, ipClient string) (map[string]string, error)
	RefreshTokens(token, ipClient string) (map[string]string, error)
	ParseRefreshToken(token string) (string, error)
	GetUserEmail(token string) (string, error)
	create(token, userId string) error
}
