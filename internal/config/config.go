package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Env         string
	DatabaseURL string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if port[0] != ':' {
		port = ":" + port
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	return &Config{
		Port:        port,
		Env:         env,
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}
