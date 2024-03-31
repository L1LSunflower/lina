package config

import (
	"log"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Workers  int `env:"WORKERS"`
	Chrome   Chrome
	Telegram Telegram
	Redis    Redis
}

type Database struct {
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
	TelegramApiKey string `env:"TELEGRAM_AP_KEY"`
	TelegramKey    string `env:"TELEGRAM_KEY"`
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
