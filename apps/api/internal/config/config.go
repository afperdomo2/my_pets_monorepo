package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Port               string
	GinMode            string
	DatabaseURL        string
	JWTSecret          string
	GoogleClientID     string
	GoogleClientSecret string
	// AppURL is the public base URL of this API server (used for OAuth redirect URIs).
	AppURL string
	// FrontendURL is the base URL of the frontend app (used for post-OAuth redirect).
	FrontendURL string
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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-in-production"
		log.Println("WARNING: JWT_SECRET not set, using insecure default. Set JWT_SECRET in production.")
	}

	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:8080"
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	return &Config{
		Port:               port,
		GinMode:            os.Getenv("GIN_MODE"),
		DatabaseURL:        dsn,
		JWTSecret:          jwtSecret,
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		AppURL:             appURL,
		FrontendURL:        frontendURL,
	}
}
