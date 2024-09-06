package service

import "JwtTestTask/internal/repository"

type Service struct {
	IAuthService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		IAuthService: NewAuthService(repo),
	}
}
