package handler

import (
	"JwtTestTask/internal/domain/models"
	"JwtTestTask/internal/services"
	"JwtTestTask/internal/services/auth"
	"JwtTestTask/internal/storage/postgres"
)

type Auth interface {
	CreateUser(user models.User) (string, error)
	GenerateTokens(userId, ipClient string) (map[string]string, error)
	SendMessageEmail(email, message string) error
	RefreshTokens(token services.Tokens, ipClient string) (map[string]string, error)
	GetUserEmail(token string) (string, error)
}

type Service struct {
	Auth
}

func NewService(storage *postgres.Storage) *Service {
	return &Service{
		Auth: auth.New(storage, storage),
	}
}
