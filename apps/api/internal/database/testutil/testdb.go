package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/my-pets/api/internal/database"
)

// NewPostgres crea un contenedor PostgreSQL 16, ejecuta AutoMigrate y retorna
// el *gorm.DB listo para tests, junto con una función de limpieza.
func NewPostgres(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("testutil.NewPostgres RunContainer: %v", err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("testutil.NewPostgres ConnectionString: %v", err)
	}

	db, err := gorm.Open(gormPostgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("testutil.NewPostgres gorm.Open: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("testutil.NewPostgres db.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(5)

	if err := database.Migrate(db); err != nil {
		t.Fatalf("testutil.NewPostgres Migrate: %v", err)
	}

	cleanup := func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("testutil.NewPostgres cleanup: %v", err)
		}
	}

	return db, cleanup
}
