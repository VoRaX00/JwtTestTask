package postgres

import "github.com/jmoiron/sqlx"

func New(storagePath string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", storagePath)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
