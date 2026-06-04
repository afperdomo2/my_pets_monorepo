# Convenciones — Backend Go

## Estructura de cada dominio

```
internal/domain/<name>/
├── repository.go    # Interface + ErrNotFound sentinel error
├── gorm_repo.go     # Implementación GORM (struct gormRepo, constructor NewGormRepo)
├── handler.go       # Handler struct con métodos para cada endpoint
├── payload.go       # Structs request con tags binding:"required"
└── routes.go        # RegisterRoutes(rg *gin.RouterGroup, h *Handler)
```

## Imports (3 bloques)

```go
import (
    "net/http"                    // 1. stdlib
    "github.com/gin-gonic/gin"    // 2. terceros
    "github.com/my-pets/api/internal/models" // 3. internos
)
```

## Structs y tags

```go
type Pet struct {
    ID   string `json:"id"   gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Name string `json:"name" gorm:"not null"`
}
```

- `json` siempre snake_case
- `binding:"required"` en payload structs
- Timestamps como `time.Time`

## Handlers

- Naming: `VerbNoun` (`GetPets`, `UpdatePet`)
- Constructor recibe `*config.Config` + repositorios de otros dominios según necesite
- Scope de reads/writes por `ownerID` (extraído del JWT)

## Errores

- Validación: `validation.Translate(err)` → mensaje en español + HTTP 400
- No encontrado: `errors.Is(err, ErrNotFound)` → HTTP 404
- `log.Fatalf` solo en `main.go`, nunca `panic()` en handlers

## Respuestas

```go
// Lista paginada:
gin.H{"data": slice, "total": N, "page": 1, "per_page": 10, "total_pages": 3}
// Recurso individual:
gin.H{"data": pet}
// Eliminación:
gin.H{"message": "registro eliminado"}
// Error:
gin.H{"error": "descripción"}
```
