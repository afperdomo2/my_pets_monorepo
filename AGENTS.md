# AGENTS.md

Guía compacta para agentes en este monorepo. Detalles por tema en `docs/`.

## Stack

| Capa | Tecnología |
|---|---|
| Backend | Go 1.25 + Gin + GORM |
| Frontend | Vue 3 + Vite 7 + TypeScript + TanStack Query |
| Monorepo | pnpm workspaces + Turborepo |
| Base de datos | Postgres 16 |
| Contenedores | Docker Compose (3 servicios: `db`, `api`, `web`) |

## Skills disponibles (`.agents/skills/`)

Cargar con `skill` tool según la tarea:

| Skill | Cuándo usarlo |
|---|---|
| `vue` | Componentes SFC, Composition API, macros, reactividad |
| `vue-best-practices` | **Obligatorio** para cualquier tarea Vue |
| `vue-debug-guides` | Debugging runtime/SSR/hidratación Vue |
| `vue-pinia-best-practices` | Stores Pinia, estado global |
| `vite` | Configuración Vite, plugins, SSR, Rolldown |
| `zod` | Schemas Zod, validación, tipos inferidos |
| `turborepo` | Pipeline tasks, turbo.json, caché, filtros |
| `typescript-advanced-types` | Tipos complejos, genéricos, condicionales |
| `frontend-design` | UI/UX, componentes visuales, diseño |
| `accessibility` | Auditoría a11y WCAG 2.2 |
| `seo` | Meta tags, datos estructurados, sitemap |
| `nodejs-best-practices` | Patrones Node.js, async, seguridad |

## Context7 MCP — documentación de librerías

Usar **siempre** para librerías nuevas o dudas de API:

1. `context7_resolve-library-id` con nombre exacto
2. `context7_query-docs` con el ID y la pregunta

## Comandos

### Backend (raíz, vía Makefile)

| Comando | Descripción | Docker |
|---|---|---|
| `make dev-api` | Hot-reload con air | ❌ |
| `make test-api` | Tests handler + unit (~0.1s) | ❌ |
| `make test-api-integration` | Handler + unit + repo (~5 min) | ✅ |
| `make lint-api` | `go vet ./...` | ❌ |
| `make swag` | Regenerar Swagger docs | ❌ |
| `make tidy` | `go mod tidy` | ❌ |
| `make build-api` | Compilar binario `bin/server` | ❌ |
| `make docker-dev` | Postgres + API + frontend | ✅ |

### Frontend (desde `apps/web/`)

| Comando | Descripción |
|---|---|
| `pnpm dev` | Vite dev server `:3000` (proxy `/api/*` → `:8080`) |
| `pnpm type-check` | `vue-tsc --build` (**obligatorio** tras cambios `.ts`/`.vue`) |
| `pnpm lint` | `run-s lint:*` (oxlint primero, eslint después) |
| `pnpm build` | `run-p type-check "build-only"` |

### Raíz

| Comando | Descripción |
|---|---|
| `pnpm dev` | Turborepo dev (frontend) |
| `pnpm build` | Turborepo build |
| `pnpm lint` | Turborepo lint |

## Reglas

- **No commits automáticos** — solo cuando el usuario lo pida
- **Documentación en español** — comentarios, docstrings, commits
- **Type-check obligatorio** tras cambios en `.ts`/`.vue` (no confiar solo en linter)
- **Swagger**: `make swag` tras crear/modificar/eliminar endpoints
- **GORM AutoMigrate** corre al iniciar (`internal/database/`)
- **authStore.initSession()** se ejecuta **antes** de `app.mount()` en `main.ts`
- Leer `docs/rules.md` — enums (const objects, no TS `enum`), paginación, lint
- Leer `docs/web/pitfalls.md` — doble request al refrescar listas, servicios sin HTTP compartido
- **Nuevas librerías**: seguir checklist en `docs/library-vetting.md` (reputación, bundle, licencia, Context7)

## Arquitectura

| Ruta | Contenido |
|---|---|
| `apps/api/cmd/server/main.go` | Entrypoint backend |
| `apps/api/internal/domain/<name>/` | Cada dominio: handler, repository (interface + GORM), payload, routes |
| `apps/api/internal/server/server.go` | Cableado: crea repos, instancia handlers, registra rutas |
| `apps/web/src/main.ts` | Entrypoint frontend |
| `apps/web/src/composables/` | TanStack Query hooks (CRUD + cache) |
| `apps/web/src/services/` | Llamadas HTTP (`get/post/put/patch/del` desde `@/services/http`) |
| `apps/web/src/schemas/` | Zod schemas para formularios (vee-validate) |
| `apps/web/src/stores/` | Pinia (sesión, UI) |

**Excepciones HTTP**: `healthCatalogService` y `setupService` tienen su propio `request()`, no usan `@/services/http`.

## Pre-commit hooks (Lefthook)

| Job | Se ejecuta si cambian |
|---|---|
| `lint-api` | `apps/api/**/*.go` |
| `test-api` | `apps/api/**/*.go` |
| `test-api-integration` | `apps/api/**/gorm_repo.go` |
| `lint-web` | `apps/web/**/*.{ts,vue}` |
| `type-check` | `apps/web/**/*.{ts,vue}` |

Los repo tests solo se ejecutan si cambió algún `gorm_repo.go`. Config: `lefthook.yml`.

## Testing

### Backend (Go)

| Tipo | Stack | Archivo | Build tag |
|---|---|---|---|
| Handler | `testify/mock` + `gin.CreateTestContext` | `handler_test.go` | — |
| Pure unit | `testing` + `testify/require` | `*_test.go` | — |
| Repository | `testcontainers-go` + Postgres 16 | `gorm_repo_test.go` | `integration` |

- Sin DB en handler tests (repos mockeados via interfaces)
- Repo tests levantan Postgres real en Docker, solo con `-tags=integration`
- Tras cambios backend: `make test-api` → `make lint-api`
- Repo tests solo bajo pedido explícito o antes de mergear
- Actualizar tabla de cobertura en `docs/api/testing.md` al crear/modificar tests
- Seed de datos: UUIDs con helper `uid(n)`, fechas con `time.Now().UTC()`

### Frontend

Vitest no configurado.

## Docker Compose

| Archivo | Incluye DB | Uso |
|---|---|---|
| `docker-compose.prod.yml` | ✅ Postgres interno | Producción local, build multi-stage |
| `docker-compose.cloud.yml` | ❌ RDS externo | AWS / Terraform |

## Convenciones clave

- **Enums**: `const` objects, NO TypeScript `enum`
- **Auth**: Cookies HttpOnly (`access_token` 20 min, `refresh_token` 20 días), no localStorage
- **Contratos FE/BE**: tipos TypeScript deben reflejar exactamente structs Go
- **Redes restringidas**: `GONOSUMDB="*" go build ./...` si hay problemas

## Datos de prueba

`FAKE_DATA.md` contiene scripts SQL y Go para crear 20 usuarios con contraseña `Password123!`.
