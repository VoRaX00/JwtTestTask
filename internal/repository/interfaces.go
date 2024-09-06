package repository

import "JwtTestTask/models"

type IAuthRepository interface {
	Get(user models.User) (string, error)
	Create(user models.User) (string, error)
	RefreshTokens(id string) (map[string]string, error)
}
