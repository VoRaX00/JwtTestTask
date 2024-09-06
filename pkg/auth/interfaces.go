package auth

import "time"

type TokenManager interface {
	NewAccessToken(ipClient string, ttl time.Duration) (string, error)
	NewRefreshToken(ipClient string, ttl time.Duration) (string, error)
}
