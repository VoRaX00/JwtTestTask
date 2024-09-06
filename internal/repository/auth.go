package repository

import (
	"JwtTestTask/models"
	"fmt"
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

func (r *AuthRepository) Get(user models.User) (string, error) {
	var id string
	query := fmt.Sprintf("SELECT id FROM Users WHERE email=$1 AND password_hash=$2")
	err := r.db.Get(&id, query, user.Email, user.Password)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *AuthRepository) Create(user models.User) (string, error) {
	query := fmt.Sprintf("INSERT INTO Users (id, email, password_hash) VALUES ($1, $2, $3)")
	err := r.db.QueryRow(query, user.Id, user.Email, user.Password)
	if err != nil {
		return "", nil
	}
	return user.Id, nil
}

func (r *AuthRepository) RefreshTokens(id string) (map[string]string, error) {
	return nil, nil
}
