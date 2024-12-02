package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env             string        `yaml:"env" env-default:"local"`
	AccessTokenTTL  time.Duration `yaml:"token_ttl" env-default:"1h"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env-default:"2h"`
	Server          CfgServer     `yaml:"server"`
	DB              CfgDB         `yaml:"db"`
}

type CfgServer struct {
	Port    int           `yaml:"port" env-default:"8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

type CfgDB struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"-"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"ssl_mode"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config file path empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("Error read config: " + err.Error())
	}

	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	if cfg.DB.Password == "" {
		panic("Password is empty")
	}
	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "config path")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
