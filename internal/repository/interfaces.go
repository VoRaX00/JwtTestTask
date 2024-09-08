package repository

import (
	"JwtTestTask/models"
)

type IAuthRepository interface {
	Get(user models.User) (string, error)
	Create(user models.User) (string, error)
	AddToken(token models.RefreshToken) error
}
