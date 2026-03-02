package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Port        string
	GinMode     string
	DatabaseURL string
}

// Load reads the .env file (if present) and returns a populated Config.
// Falls back to sensible defaults for local development.
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading environment variables directly")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=devuser password=devpassword dbname=my_pets port=5432 sslmode=disable"
	}

	return &Config{
		Port:        port,
		GinMode:     os.Getenv("GIN_MODE"),
		DatabaseURL: dsn,
	}
}
