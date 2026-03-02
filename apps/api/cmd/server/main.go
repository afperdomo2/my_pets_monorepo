// Package main is the entrypoint for the My Pets API server.
package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/my-pets/api/internal/config"
	"github.com/my-pets/api/internal/server"
)

//	@title			My Pets API
//	@version		1.0
//	@description	API REST para gestión de mascotas y usuarios.

//	@host		localhost:8080
//	@BasePath	/

//	@tag.name			health
//	@tag.description	Liveness check
//	@tag.name			pets
//	@tag.description	Operaciones sobre mascotas
//	@tag.name			users
//	@tag.description	Operaciones sobre usuarios

//	@securityDefinitions.apikey	CookieAuth
//	@in							cookie
//	@name						access_token
//	@description				Cookie HttpOnly de sesión. Obtenida al hacer login o refresh. Se envía automáticamente por el navegador.

func main() {
	cfg := config.Load()

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	server.Run(cfg, db)
}
