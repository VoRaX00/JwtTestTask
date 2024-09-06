package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserIp string `json:"user_ip"`
}

type Manager struct {
	signInKey string
}

func NewManager(signInKey string) (*Manager, error) {
	if signInKey == "" {
		return nil, errors.New("signInKey is empty")
	}
	return &Manager{signInKey: signInKey}, nil
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
	return token.SignedString([]byte(m.signInKey))
}

func (m *Manager) NewRefreshToken(ipClient string, ttl time.Duration) (string, error) {

}
