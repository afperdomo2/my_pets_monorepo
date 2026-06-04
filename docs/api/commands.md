# Comandos — Backend Go

Todos los comandos se ejecutan desde la **raíz del monorepo** (vía `Makefile`).

## Desarrollo

| Comando | Descripción |
|---|---|
| `make dev-api` | Hot-reload con air |
| `make build-api` | Compilar binario → `bin/server` |
| `make test-api` | `go test ./... -v` |
| `make lint-api` | `go vet ./...` |
| `make swag` | Regenerar Swagger docs |
| `make tidy` | `go mod tidy` |

## Tests

```bash
# Test enfocado en un dominio
go test ./internal/domain/<name>/ -run <TestName> -v

# Todos los tests
make test-api
```

> **No existen tests Go** en el proyecto — no hay archivos `_test.go`.

## Build con restricciones de red

```bash
GONOSUMDB="*" go build ./...
```
