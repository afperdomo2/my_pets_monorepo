# Comandos — Backend Go

Todos los comandos se ejecutan desde la **raíz del monorepo** (vía `Makefile`).

## Desarrollo

| Comando | Descripción |
|---|---|
| `make dev-api` | Hot-reload con air |
| `make build-api` | Compilar binario → `bin/server` |
| `make test-api` | `go test ./... -v` |
| `make lint-api` | `go vet ./...` |
| `make swag` | Regenerar Swagger docs (`apps/api/docs/`) |
| `make tidy` | `go mod tidy` |

## Tests

```bash
# Test enfocado en un dominio
go test ./internal/domain/<name>/ -run <TestName> -v

# Todos los tests
make test-api
```

### Test patterns usados

| Tipo | Descripción | Ubicación |
|---|---|---|
| Handler tests | Mocks de repositorios con `testify/mock` + `gin.CreateTestContext` | `internal/domain/<name>/handler_test.go` |
| Pure unit tests | Tests sin mocks ni DB (cálculos, validaciones) | `internal/domain/<name>/xxx_test.go` |

### Stack de testing

- [`testify`](https://github.com/stretchr/testify) — mocks (`testify/mock`) y aserciones (`testify/require`)
- `httptest` — `ResponseRecorder` + `gin.CreateTestContext`
- Sin DB en tests de handlers (repositorios mockeados)

## Swagger / Documentación de endpoints

- Cada vez que se **cree, modifique o elimine un endpoint**, ejecutar `make swag` para regenerar los docs.
- Los handlers usan anotaciones `// @Summary`, `// @Param`, `// @Success`, `// @Router` (ver ejemplos en los handlers existentes).
- Los archivos generados están en `apps/api/docs/` (`docs.go`, `swagger.json`, `swagger.yaml`).

## Build con restricciones de red

```bash
GONOSUMDB="*" go build ./...
```
