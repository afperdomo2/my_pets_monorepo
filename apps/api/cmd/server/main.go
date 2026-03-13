// Package main is the entrypoint for the My Pets API server.
package main

import (
	"log"

	"github.com/my-pets/api/internal/config"
	"github.com/my-pets/api/internal/database"
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
//	@tag.name			vaccines-catalog
//	@tag.description	Operaciones sobre el catálogo de vacunas

//	@securityDefinitions.apikey	CookieAuth
//	@in							cookie
//	@name						access_token
//	@description				Cookie HttpOnly de sesión. Obtenida al hacer login o refresh. Se envía automáticamente por el navegador.

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL, cfg.GinMode)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	server.Run(cfg, db)
}
