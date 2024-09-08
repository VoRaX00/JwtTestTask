package service

import (
	"JwtTestTask/internal/repository"
	"JwtTestTask/models"
	"JwtTestTask/pkg/auth"
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const (
	salt            = "ouwwlcq]3.djc.4iolor001mcrufn"
	signingKey      = "iydjiadiopejo223jn2nvuernveia.xQ!eij"
	accessTokenTTL  = time.Hour * 24
	refreshTokenTTL = time.Hour * 48
)

type AuthService struct {
	repo repository.IAuthRepository
}

func NewAuthService(repo repository.IAuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) GenerateTokens(user models.User, ipClient string) (map[string]string, error) {
	user.Password = s.generatePasswordHash(user.Password)
	_, err := s.repo.Get(user)
	if err != nil {
		return nil, err
	}

	tokens := map[string]string{}
	tokenManager, err := auth.NewManager(signingKey)
	if err != nil {
		return nil, err
	}

	accessToken, err := tokenManager.NewAccessToken(ipClient, accessTokenTTL)
	if err != nil {
		return nil, err
	}

	tokens["access_token"] = accessToken
	refreshToken, err := tokenManager.NewRefreshToken(ipClient, refreshTokenTTL)
	if err != nil {
		return nil, err
	}

	tokens["refresh_token"] = refreshToken

	err = s.addToken(refreshToken, user.Id)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (s *AuthService) Create(user models.User) (string, error) {
	user.Id = uuid.New().String()
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.Create(user)
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) addToken(token, userId string) error {
	tokenManager, err := auth.NewManager(signingKey)
	if err != nil {
		return err
	}

	hash, err := tokenManager.HashRefreshToken(token)
	if err != nil {
		return err
	}

	refreshToken := models.RefreshToken{
		UserId:    userId,
		TokenHash: hash,
		TTL:       refreshTokenTTL,
	}
	return s.repo.AddToken(refreshToken)
}
