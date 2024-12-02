package models

import "time"

type RefreshToken struct {
	Id               int       `db:"id"`
	UserId           string    `db:"user_id"`
	RefreshTokenHash string    `db:"refresh_token_hash"`
	Ip               string    `db:"ip"`
	CreatedAt        time.Time `db:"created_at"`
	ExpiresAt        time.Time `db:"expires_at"`
}
