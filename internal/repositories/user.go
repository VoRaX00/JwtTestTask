package repositories

import (
	"JwtTestTask/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type UsersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}

func (r *UsersRepository) Create(user models.User) (string, error) {
	query := fmt.Sprintf("INSERT INTO Users (id, email, password_hash) VALUES ($1, $2, $3) RETURNING id")
	row := r.db.QueryRow(query, user.Id, user.Email, user.Password)
	if row.Scan(&user.Id) != nil {
		return "", nil
	}
	return user.Id, nil
}
