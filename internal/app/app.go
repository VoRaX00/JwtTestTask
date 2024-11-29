package app

import (
	"JwtTestTask/internal/server"
	"JwtTestTask/internal/storage/postgres"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
	"log/slog"
)

type App struct {
	Server *server.Server
	DB     *sqlx.DB
}

func New(log *slog.Logger, storagePath, port string) *App {
	//db, err := connectDB(storagePath)
	//if err != nil {
	//	log.Warn(err.Error())
	//}

	//repos :=
	return nil
}

func connectDB(storagePath string) (*sqlx.DB, error) {
	db, err := postgres.New(storagePath)
	if err != nil {
		panic("error connecting to database: " + err.Error())
	}

	if err = goose.Up(db.DB, "./storage/migrations"); err != nil {
		return db, fmt.Errorf("error upgrading database: %v", err)
	}
	return db, nil
}
