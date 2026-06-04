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

> **No existen tests Go** en el proyecto — no hay archivos `_test.go`.

## Swagger / Documentación de endpoints

- Cada vez que se **cree, modifique o elimine un endpoint**, ejecutar `make swag` para regenerar los docs.
- Los handlers usan anotaciones `// @Summary`, `// @Param`, `// @Success`, `// @Router` (ver ejemplos en los handlers existentes).
- Los archivos generados están en `apps/api/docs/` (`docs.go`, `swagger.json`, `swagger.yaml`).

## Build con restricciones de red

```bash
GONOSUMDB="*" go build ./...
```
