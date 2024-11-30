package handler

import (
	"JwtTestTask/internal/services/auth"
	"JwtTestTask/internal/storage"
)

type Service struct {
	auth.IUserService
	auth.ITokenService
}

func NewService(repo *storage.Repository) *Service {
	return &Service{
		ITokenService: auth.NewTokenService(repo.Token),
		IUserService:  auth.NewUserService(repo.User),
	}
}
