package app

import (
	"JwtTestTask/internal/app/server"
	"JwtTestTask/internal/config"
	"JwtTestTask/internal/handler"
	"JwtTestTask/internal/storage"
	"JwtTestTask/internal/storage/postgres"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
	"log/slog"
)

type App struct {
	Server *server.Server
	DB     *sqlx.DB
	log    *slog.Logger
}

func New(log *slog.Logger, storagePath string, cfg config.CfgServer) *App {
	db, err := connectDB(storagePath)
	if err != nil {
		log.Warn(err.Error())
	}

	repos := storage.NewRepository(db)
	service := handler.NewService(repos)
	handlers := handler.NewHandler(service)
	srv := server.New(log, cfg, handlers.InitRoutes())
	return &App{
		Server: srv,
		DB:     db,
		log:    log,
	}
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

func (a *App) Stop(ctx context.Context) {
	err := a.DB.Close()
	if err != nil {
		a.log.Warn(err.Error())
	}
	a.Server.Stop(ctx)
}
