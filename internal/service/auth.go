package service

import (
	"JwtTestTask/internal/repository"
	"JwtTestTask/models"
)

type AuthService struct {
	repo repository.IAuthRepository
}

func NewAuthService(repo repository.IAuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) SignIn(id string) (map[string]interface{}, error) {
	return s.repo.SignIn(id)
}

func (s *AuthService) SignUp(user models.User) (string, error) {
	return s.repo.SignUp(user)
}

func (s *AuthService) RefreshToken(id string) (map[string]interface{}, error) {
	return s.repo.RefreshToken(id)
}
