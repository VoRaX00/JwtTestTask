package repository

import (
	"JwtTestTask/models"
	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) SignIn(id string) (map[string]interface{}, error) {
	return nil, nil
}

func (r *AuthRepository) SignUp(user models.User) (string, error) {
	return "", nil
}

func (r *AuthRepository) RefreshToken(id string) (map[string]interface{}, error) {
	return nil, nil
}
