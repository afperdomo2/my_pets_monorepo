# Comandos — Backend Go

Todos los comandos se ejecutan desde la **raíz del monorepo** (vía `Makefile`).

## Desarrollo

| Comando | Descripción | Docker |
|---|---|---|---|
| `make dev-api` | Hot-reload con air | ❌ |
| `make build-api` | Compilar binario → `bin/server` | ❌ |
| `make test-api` | Handler + unit tests (~0.1s) | ❌ |
| `make test-api-integration` | Handler + unit + repo tests (~5 min) | ✅ |
| `make lint-api` | `go vet ./...` | ❌ |
| `make swag` | Regenerar Swagger docs (`apps/api/docs/`) | ❌ |
| `make tidy` | `go mod tidy` | ❌ |

## Tests

Ver documentación detallada en [`docs/api/testing.md`](testing.md).

```bash
# Handler + unit (rápido, sin Docker)
make test-api

# Handler + unit + repo (requiere Docker)
make test-api-integration

# Test enfocado en un dominio
go test ./internal/domain/<name>/ -run <TestName> -v

# Test de repositorio (build tag integration)
go test -tags=integration ./internal/domain/<name>/ -run TestGormRepo -v
```

### Test patterns usados

| Tipo | Descripción | Ubicación | Build tag |
|---|---|---|---|---|
| Handler tests | Mocks de repositorios con `testify/mock` + `gin.CreateTestContext` | `internal/domain/<name>/handler_test.go` | — |
| Pure unit tests | Tests sin mocks ni DB (cálculos, validaciones) | `internal/domain/<name>/xxx_test.go` | — |
| Repository tests | Queries reales contra Postgres via testcontainers | `internal/domain/<name>/gorm_repo_test.go` | `integration` |

### Stack de testing

- [`testify`](https://github.com/stretchr/testify) — mocks (`testify/mock`) y aserciones (`testify/require`)
- `httptest` — `ResponseRecorder` + `gin.CreateTestContext`
- Sin DB en tests de handlers (repositorios mockeados)
- [`testcontainers-go`](https://golang.testcontainers.org/) — Postgres 16 en contenedor para tests de repositorio
- Los tests de repositorio se excluyen por defecto con `//go:build integration`; se ejecutan con `-tags=integration`

## Swagger / Documentación de endpoints

- Cada vez que se **cree, modifique o elimine un endpoint**, ejecutar `make swag` para regenerar los docs.
- Los handlers usan anotaciones `// @Summary`, `// @Param`, `// @Success`, `// @Router` (ver ejemplos en los handlers existentes).
- Los archivos generados están en `apps/api/docs/` (`docs.go`, `swagger.json`, `swagger.yaml`).

## Build con restricciones de red

```bash
GONOSUMDB="*" go build ./...
```
