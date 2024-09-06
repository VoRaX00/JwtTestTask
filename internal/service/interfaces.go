package service

import "JwtTestTask/models"

type IAuthService interface {
	SignIn(id string) (map[string]interface{}, error)
	SignUp(user models.User) (string, error)
	RefreshToken(id string) (map[string]interface{}, error)
}
