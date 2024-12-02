package auth

import (
	"JwtTestTask/internal/domain/models"
	"JwtTestTask/internal/services"
	"JwtTestTask/pkg/manager"
	"crypto/sha256"
	"database/sql"
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
	InvalidToken  = errors.New("invalid token")
	IpNotEquals   = errors.New("ip not equal")
	TokenExpired  = errors.New("token expired")
	TokenNotFound = errors.New("token not found")
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

func (s *Auth) GenerateTokens(userId, ipClient string) (services.Tokens, error) {
	var tokens services.Tokens
	accessToken, err := s.tokenManager.NewAccessToken(ipClient, accessTokenTTL)
	if err != nil {
		return services.Tokens{}, err
	}

	tokens.AccessToken = accessToken
	refreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return services.Tokens{}, err
	}

	tokens.RefreshToken = refreshToken

	hash, err := s.tokenManager.HashRefreshToken(refreshToken)
	foundToken, err := s.tokenProvider.GetRefreshTokenByUserId(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = s.create(hash, userId, ipClient)
			if err != nil {
				return services.Tokens{}, err
			}
		} else {
			return services.Tokens{}, err
		}
	}

	err = s.tokenProvider.RefreshToken(hash, foundToken.RefreshTokenHash, ipClient, refreshTokenTTL)
	if err != nil {
		return services.Tokens{}, err
	}
	return tokens, nil
}

func (s *Auth) create(tokenHash, userId, ipClient string) error {
	expiresAt := time.Now().Add(refreshTokenTTL)
	refreshToken := models.RefreshToken{
		UserId:           userId,
		RefreshTokenHash: tokenHash,
		Ip:               ipClient,
		ExpiresAt:        expiresAt,
	}
	return s.tokenProvider.CreateToken(refreshToken)
}

func (s *Auth) RefreshTokens(tokens services.Tokens, ipClient string) (services.Tokens, error) {
	const op = "auth.RefreshTokens"

	err := s.validateAccessToken(tokens.AccessToken, ipClient)
	if err != nil {
		return services.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	hashToken, err := s.tokenManager.HashRefreshToken(tokens.RefreshToken)
	if err != nil {
		return services.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	err = s.validateRefreshToken(hashToken, ipClient)
	if err != nil {
		return services.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	refreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return services.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	hashNewToken, err := s.tokenManager.HashRefreshToken(refreshToken)
	if err != nil {
		return services.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	err = s.tokenProvider.RefreshToken(hashNewToken, hashToken, ipClient, refreshTokenTTL)
	if err != nil {
		return services.Tokens{}, err
	}

	accessToken, err := s.tokenManager.NewAccessToken(ipClient, accessTokenTTL)
	if err != nil {
		return services.Tokens{}, err
	}

	return services.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Auth) validateAccessToken(token string, ipClient string) error {
	claims, err := s.tokenManager.DecodeAccessToken(token)
	if err != nil {
		return err
	}

	if claims == nil {
		return InvalidToken
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return TokenExpired
	}

	if ipClient != claims.UserIp {
		return IpNotEquals
	}
	return nil
}

const emailWarning = "В ваш аккаунт зашли с другого устройства"

func (s *Auth) validateRefreshToken(hashToken, ipClient string) error {
	token, err := s.tokenProvider.GetRefreshToken(hashToken)
	if err != nil {
		return err
	}

	if token == nil {
		return TokenNotFound
	}

	if token.Ip != ipClient {
		email, err := s.tokenProvider.GetUserEmail(token.RefreshTokenHash)
		if err != nil {
			return err
		}

		err = s.SendMessageEmail(email, emailWarning)
		if err != nil {
			return err
		}
		return IpNotEquals
	}

	if time.Now().After(token.ExpiresAt) {
		return TokenExpired
	}
	return nil
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
