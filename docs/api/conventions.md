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

## Tests

### Handler tests (`handler_test.go`)

- Usar `testify/mock` para mockear interfaces de repositorio
- Usar `gin.CreateTestContext(w)` + `httptest.NewRecorder()` para simular requests HTTP
- Definir `mockRepo` anónimo dentro del mismo package (acceso a tipos no exportados)
- Patrón: `setup → act → assert` con `mock.AssertExpectations(t)`
- Test name en snake_case con el patrón `Test<HandlerName>_<Scenario>`
- NO usar bases de datos reales — los repos se mockean siempre en handler tests

### Pure unit tests (`xxx_test.go`)

- Tests sin mocks ni DB para lógica pura (cálculos, validaciones, helpers)
- Usar `testify/require` para aserciones
- Preferir table-driven tests con `t.Run()` por subtest

```go
func TestCalculateDogLifeStage(t *testing.T) {
    tests := []struct {
        name     string
        birth    time.Time
        size     SizeCategory
        expected string
    }{...}
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            require.Equal(t, tt.expected, CalculateDogLifeStage(tt.birth, tt.size))
        })
    }
}
```

### Stack

| Herramienta | Propósito |
|---|---|
| `testing` | Runner estándar de Go |
| `testify/mock` | Creación de mocks |
| `testify/require` | Aserciones con fail inmediato |
| `net/http/httptest` | `ResponseRecorder` para capturar respuestas HTTP |
| `gin.CreateTestContext` | Contexto Gin sin servidor HTTP real |

## Swagger

- Todo endpoint debe tener anotaciones Swagger (`@Summary`, `@Tags`, `@Param`, `@Success`, `@Router`) en su handler.
- Después de crear/modificar/eliminar un endpoint, regenerar docs con `make swag`.

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
