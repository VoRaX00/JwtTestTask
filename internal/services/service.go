package services

import (
	"JwtTestTask/internal/repositories"
	"JwtTestTask/pkg/auth"
)

type Service struct {
	IUserService
	ITokenService
}

func NewService(repo *repositories.Repository, tokenManager auth.TokenManager) *Service {
	return &Service{
		ITokenService: NewTokenService(repo.ITokenRepository, tokenManager),
		IUserService:  NewUserService(repo.IUserRepository),
	}
}
