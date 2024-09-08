package models

import "time"

type RefreshToken struct {
	UserId    string
	TokenHash string
	TTL       time.Duration
}
