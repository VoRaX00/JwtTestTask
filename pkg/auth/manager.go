package auth

import (
	"encoding/base64"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
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

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	if _, err := r.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (m *Manager) HashRefreshToken(token string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
