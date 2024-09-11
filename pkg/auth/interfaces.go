package auth

import "time"

type TokenManager interface {
	NewAccessToken(ipClient string, ttl time.Duration) (string, error)
	NewRefreshToken() (string, error)
	HashRefreshToken(token string) (string, error)
}
