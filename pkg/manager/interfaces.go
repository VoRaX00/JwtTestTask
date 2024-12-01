package manager

import "time"

type TokenManager interface {
	NewAccessToken(ipClient string, ttl time.Duration) (string, error)
	NewRefreshToken() (string, error)
	DecodeAccessToken(accessToken string) (*TokenClaims, error)
	HashRefreshToken(token string) (string, error)
}
