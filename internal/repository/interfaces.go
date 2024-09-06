package repository

import "JwtTestTask/models"

type IAuthRepository interface {
	SignIn(id string) (map[string]interface{}, error)
	SignUp(user models.User) (string, error)
	RefreshToken(id string) (map[string]interface{}, error)
}
