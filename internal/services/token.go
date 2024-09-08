package services

import (
	"JwtTestTask/internal/repositories"
	"JwtTestTask/models"
	"JwtTestTask/pkg/auth"
	"time"
)

const (
	salt            = "ouwwlcq]3.djc.4iolor001mcrufn"
	accessTokenTTL  = time.Hour * 24
	refreshTokenTTL = time.Hour * 48
)

type TokenService struct {
	repo         repositories.ITokenRepository
	tokenManager auth.TokenManager
}

func NewTokenService(repo repositories.ITokenRepository, tokenManager auth.TokenManager) *TokenService {
	return &TokenService{
		repo:         repo,
		tokenManager: tokenManager,
	}
}

func (s *TokenService) GenerateTokens(userId, ipClient string) (map[string]string, error) {
	tokens := map[string]string{}

	accessToken, err := s.tokenManager.NewAccessToken(ipClient, accessTokenTTL)
	if err != nil {
		return nil, err
	}

	tokens["access_token"] = accessToken
	refreshToken, err := s.tokenManager.NewRefreshToken(ipClient, refreshTokenTTL)
	if err != nil {
		return nil, err
	}

	tokens["refresh_token"] = refreshToken

	err = s.create(refreshToken, userId)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (s *TokenService) ParseRefreshToken(token string) (string, error) {
	ip, err := s.tokenManager.ParseRefreshToken(token)
	return ip, err
}

func (s *TokenService) RefreshTokens(token, ipClient string) (map[string]string, error) {
	return nil, nil
}

func (s *TokenService) create(token, userId string) error {
	hash, err := s.tokenManager.HashRefreshToken(token)
	if err != nil {
		return err
	}

	refreshToken := models.RefreshToken{
		UserId:           userId,
		RefreshTokenHash: hash,
		ExpiresAt:        refreshTokenTTL,
	}
	return s.repo.Create(refreshToken)
}

func (s *TokenService) GetUserEmail(token string) (string, error) {
	hash, err := s.tokenManager.HashRefreshToken(token)
	if err != nil {
		return "", err
	}
	return s.repo.GetUserEmail(hash)
}
