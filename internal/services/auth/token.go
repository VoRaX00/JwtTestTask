package auth

import (
	"JwtTestTask/internal/services"
	authrepo "JwtTestTask/internal/storage/auth"
	"JwtTestTask/models"
	"JwtTestTask/pkg/manager"
	"net/smtp"
	"os"
	"time"
)

const (
	accessTokenTTL  = time.Hour * 24
	refreshTokenTTL = time.Hour * 48
)

type TokenService struct {
	repo         authrepo.ITokenRepository
	tokenManager manager.TokenManager
}

func NewTokenService(repo authrepo.ITokenRepository) *TokenService {
	signingKey := os.Getenv("JWT_SIGNING_KEY")
	return &TokenService{
		repo:         repo,
		tokenManager: manager.NewManager(signingKey),
	}
}

func (s *TokenService) GenerateTokens(userId, ipClient string) (map[string]string, error) {
	tokens := map[string]string{}

	accessToken, err := s.tokenManager.NewAccessToken(ipClient, accessTokenTTL)
	if err != nil {
		return nil, err
	}

	tokens["access_token"] = accessToken
	refreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, err
	}

	tokens["refresh_token"] = refreshToken
	err = s.create(refreshToken, userId, ipClient)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

const emailWarning = "В ваш аккаунт зашли с другого устройства"

func (s *TokenService) RefreshTokens(token services.Tokens, ipClient string) (map[string]string, error) {
	tokens := map[string]string{}

	refreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, err
	}

	hashNewToken, err := s.tokenManager.HashRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	hashToken, err := s.tokenManager.HashRefreshToken(token.RefreshToken)
	if err != nil {
		return nil, err
	}

	ip, err := s.repo.RefreshTokens(hashNewToken, hashToken, ipClient, refreshTokenTTL)
	if err != nil {
		return nil, err
	}

	if ip != "" {
		email, err := s.repo.GetUserEmail(hashNewToken)
		if err != nil {
			return nil, err
		}

		err = s.SendMessageEmail(email, emailWarning)
		if err != nil {
			return nil, err
		}
	}

	accessToken, err := s.tokenManager.NewAccessToken(ipClient, accessTokenTTL)
	if err != nil {
		return nil, err
	}

	tokens["access_token"] = accessToken
	tokens["refresh_token"] = refreshToken
	return tokens, nil
}

func (s *TokenService) create(token, userId, ipClient string) error {
	hash, err := s.tokenManager.HashRefreshToken(token)
	expiresAt := time.Now().Add(refreshTokenTTL)
	if err != nil {
		return err
	}

	refreshToken := models.RefreshToken{
		UserId:           userId,
		RefreshTokenHash: hash,
		Ip:               ipClient,
		ExpiresAt:        expiresAt,
	}
	return s.repo.Create(refreshToken)
}

func (s *TokenService) SendMessageEmail(email, message string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	to := []string{email}
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	byteMessage := []byte(message)
	authSmtp := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, authSmtp, from, to, byteMessage)
	return err
}

func (s *TokenService) GetUserEmail(token string) (string, error) {
	hash, err := s.tokenManager.HashRefreshToken(token)
	if err != nil {
		return "", err
	}
	return s.repo.GetUserEmail(hash)
}
