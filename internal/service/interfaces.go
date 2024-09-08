package service

import "JwtTestTask/models"

type IAuthService interface {
	GenerateTokens(user models.User, ipClient string) (map[string]string, error)
	Create(user models.User) (string, error)
	RefreshTokens(id string) (map[string]string, error)
	addToken(token models.RefreshToken) error
}
