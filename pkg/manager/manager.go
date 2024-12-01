package manager

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"time"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserIp string `json:"user_ip"`
}

type Manager struct {
	signingKey string
}

func NewManager(signInKey string) *Manager {
	return &Manager{signingKey: signInKey}
}

func (m *Manager) NewAccessToken(ipClient string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		&TokenClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(ttl).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
			UserIp: ipClient,
		},
	)
	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) DecodeAccessToken(accessToken string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	if _, err := r.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (m *Manager) HashRefreshToken(token string) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(token))
	return hex.EncodeToString(hash.Sum(nil)), nil
}
