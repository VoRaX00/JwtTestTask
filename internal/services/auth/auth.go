package auth

import (
	"JwtTestTask/internal/services"
	"JwtTestTask/models"
	"JwtTestTask/pkg/manager"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/smtp"
	"os"
	"time"
)

const (
	accessTokenTTL  = time.Hour * 24
	refreshTokenTTL = time.Hour * 48
)

var (
	IpNotEquals  = errors.New("ip not equal")
	TokenExpired = errors.New("token expired")
)

type Auth struct {
	userProvider  UserProvider
	tokenProvider TokenProvider
	tokenManager  manager.TokenManager
}

func New(userProvider UserProvider, tokenProvider TokenProvider) *Auth {
	signingKey := os.Getenv("JWT_SIGNING_KEY")
	return &Auth{
		userProvider:  userProvider,
		tokenProvider: tokenProvider,
		tokenManager:  manager.NewManager(signingKey),
	}
}

func (s *Auth) CreateUser(user models.User) (string, error) {
	user.Id = uuid.New().String()
	user.Password = s.generatePasswordHash(user.Password)
	return s.userProvider.CreateUser(user)
}

func (s *Auth) generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("SALT"))))
}

func (s *Auth) GenerateTokens(userId, ipClient string) (map[string]string, error) {
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

func (s *Auth) RefreshTokens(token services.Tokens, ipClient string) (map[string]string, error) {
	// TODO: доделать метод => исправить метод в репозитории refreshTokens
	const op = "auth.RefreshTokens"
	tokens := map[string]string{}

	claims, err := s.tokenManager.DecodeAccessToken(token.AccessToken)
	if err != nil {
		return nil, err
	}

	userIp := claims.UserIp

	hashToken, err := s.tokenManager.HashRefreshToken(token.RefreshToken)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenProvider.RefreshToken(hashToken)
	if err != nil {
		return nil, err
	}

	if refreshToken.Ip != userIp {
		email, err := s.tokenProvider.GetUserEmail(refreshToken.RefreshTokenHash)
		if err != nil {
			return nil, err
		}

		err = s.SendMessageEmail(email, emailWarning)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s: %w", op, IpNotEquals)
	}

	if !time.Now().After(refreshToken.ExpiresAt) {
		return nil, TokenExpired
	}

	newToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, err
	}

	hashNewToken, err := s.tokenManager.HashRefreshToken(newToken)
	if err != nil {
		return nil, err
	}

	_, err = s.tokenProvider.RefreshTokens(hashNewToken, hashToken, ipClient, refreshTokenTTL)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.tokenManager.NewAccessToken(ipClient, accessTokenTTL)
	if err != nil {
		return nil, err
	}

	tokens["access_token"] = accessToken
	tokens["refresh_token"] = newToken
	return tokens, nil
}

func (s *Auth) create(token, userId, ipClient string) error {
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
	return s.tokenProvider.CreateToken(refreshToken)
}

func (s *Auth) SendMessageEmail(email, message string) error {
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

func (s *Auth) GetUserEmail(token string) (string, error) {
	hash, err := s.tokenManager.HashRefreshToken(token)
	if err != nil {
		return "", err
	}
	return s.tokenProvider.GetUserEmail(hash)
}
