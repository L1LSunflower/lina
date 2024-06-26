package config

import (
	"log"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Workers  int `env:"WORKERS"`
	Database Database
	Chrome   Chrome
	Telegram Telegram
	Redis    Redis
}

type Database struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
	DBName   string `env:"DB_NAME"`
	SSLMode  bool   `env:"DB_SSL"`
}

type Chrome struct {
	Headless  bool `env:"HEADLESS"`
	DebugMode bool `env:"DEBUG_MODE"`
}

type Telegram struct {
	TelegramUrl string `env:"TELEGRAM_URL"`
	TelegramKey string `env:"TELEGRAM_KEY"`
}

type Redis struct {
	Host     string `env:"REDIS_HOST"`
	Password string `env:"REDIS_PASSWORD"`
	Tls      bool   `env:"REDIS_TLS"`
}

func Conf() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("FAILED TO PARSE ENV WITH ERROR: %s", err)
	}
	return cfg
}
