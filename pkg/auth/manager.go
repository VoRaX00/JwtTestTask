package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserIp string `json:"user_ip"`
}

type Manager struct {
	signingKey string
}

func NewManager(signInKey string) (*Manager, error) {
	if signInKey == "" {
		return nil, errors.New("signInKey is empty")
	}
	return &Manager{signingKey: signInKey}, nil
}

func (m *Manager) NewAccessToken(ipClient string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		&tokenClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(ttl).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
			UserIp: ipClient,
		},
	)
	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) NewRefreshToken(ipClient string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodES256, &tokenClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(ttl).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
			UserIp: ipClient,
		})
	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) HashRefreshToken(token string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
