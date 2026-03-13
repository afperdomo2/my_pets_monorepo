package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/my-pets/api/internal/models"
)

// Connect opens a GORM connection to PostgreSQL using the given DSN.
func Connect(dsn string, ginMode string) (*gorm.DB, error) {
	logLevel := logger.Silent
	if ginMode != "release" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("database.Connect: %w", err)
	}

	return db, nil
}

// Migrate runs AutoMigrate for all registered models.
// It is idempotent: safe to call on every application start.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.User{},
		&models.Pet{},
		&models.HealthCatalog{},
	); err != nil {
		return fmt.Errorf("database.Migrate: %w", err)
	}
	return nil
}
