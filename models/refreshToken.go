package models

import "time"

type RefreshToken struct {
	Id               int           `db:"id"`
	UserId           string        `db:"user_id"`
	RefreshTokenHash string        `db:"refresh_token_hash"`
	CreatedAt        time.Time     `db:"created_at"`
	ExpiresAt        time.Duration `db:"expires_at"`
	Revoked          bool          `db:"revoked"`
}
