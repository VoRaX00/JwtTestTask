package main

import (
	"JwtTestTask/internal/handler"
	"JwtTestTask/internal/repositories"
	"JwtTestTask/internal/server"
	"JwtTestTask/internal/services"
	"JwtTestTask/pkg/auth"
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("Init config error, %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Load env error, %s", err.Error())
	}

	cfg := repositories.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}

	db, err := repositories.NewPostgresDB(cfg)
	if err != nil {
		logrus.Fatalf("Init db error, %s", err.Error())
	}

	repos := repositories.NewRepository(db)

	signingKey := os.Getenv("JWT_SIGNING_KEY")
	tokenManager, err := auth.NewManager(signingKey)
	if err != nil {
		logrus.Fatalf("Init token manager error, %s", err.Error())
	}
	services := services.NewService(repos, tokenManager)

	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	go func() {
		if err = srv.Run(viper.GetString("server.port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("Run server error, %s", err.Error())
		}
	}()

	logrus.Info("Server init")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutdown server ...")
	if err = srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("Server shutdown error, %s", err.Error())
	}

	if err = db.Close(); err != nil {
		logrus.Fatalf("Database close error, %s", err.Error())
	}
}
