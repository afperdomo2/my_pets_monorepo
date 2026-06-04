# AGENTS.md

Guía compacta para agentes en este monorepo. El detalle por tema está en `docs/`.

## Stack

| Capa | Tecnología |
|---|---|
| Backend | Go 1.25 + Gin + GORM |
| Frontend | Vue 3 + Vite 7 + TypeScript + TanStack Query |
| Monorepo | pnpm workspaces + Turborepo |
| Base de datos | Postgres 16 |
| Contenedores | Docker Compose (3 servicios: `db`, `api`, `web`) |

## Skills disponibles (`.agents/skills/`)

Carga el skill con `skill` tool según la tarea:

| Skill | Cuándo usarlo |
|---|---|
| `vue` | Componentes SFC, Composition API, macros, reactividad |
| `vue-best-practices` | Cualquier tarea Vue (obligatorio leerlo) |
| `vue-debug-guides` | Debugging de errores runtime/SSR/hidratación Vue |
| `vue-pinia-best-practices` | Stores Pinia, estado global |
| `vite` | Configuración de Vite, plugins, SSR, Rolldown |
| `zod` | Schemas Zod, validación, tipos inferidos |
| `turborepo` | Pipeline tasks, turbo.json, caché, filtros |
| `typescript-advanced-types` | Tipos complejos, genéricos, condicionales |
| `frontend-design` | UI/UX, componentes visuales, diseño |
| `accessibility` | Auditoría a11y WCAG 2.2 |
| `seo` | Meta tags, datos estructurados, sitemap |
| `nodejs-best-practices` | Patrones Node.js, async, seguridad |

## Context7 MCP — documentación de librerías

Usar **siempre** que se necesite instalar, consultar o investigar librerías nuevas:

1. `context7_resolve-library-id` con el nombre exacto y el contexto de uso
2. `context7_query-docs` para obtener documentación actualizada

Esto evita depender de datos de entrenamiento que pueden estar desactualizados. Aplica a cualquier librería (Vue, Zod, TanStack Query, Pinia, Vite, etc.).

## Comandos rápidos

### Backend (desde raíz, vía Makefile)

| Comando | Descripción |
|---|---|
| `make dev-api` | Hot-reload con air |
| `make test-api` | `go test ./... -v` |
| `make lint-api` | `go vet ./...` |
| `make swag` | Regenerar Swagger docs (tras cambios en handlers) |
| `make tidy` | `go mod tidy` |
| `make build-api` | Compilar binario `bin/server` |
| `make docker-dev` | Levantar todo (Postgres + API + frontend) |

### Frontend (desde `apps/web/`)

| Comando | Descripción |
|---|---|
| `pnpm dev` | Vite dev server `:3000` |
| `pnpm type-check` | `vue-tsc --build` (**obligatorio** tras cambios `.ts`/`.vue`) |
| `pnpm lint:oxlint` | Ejecutar **primero** |
| `pnpm lint:eslint` | Ejecutar **después** |
| `pnpm build` | `run-p type-check "build-only"` |

### Raíz del monorepo

| Comando | Descripción |
|---|---|
| `pnpm dev` | Turborepo dev |
| `pnpm build` | Turborepo build |
| `pnpm lint` | Turborepo lint |

## Reglas importantes

### Generales

- **No commits automáticos** — solo cuando el usuario lo pida explícitamente
- **Documentación en español** — comentarios, docstrings, commits
- Leer `docs/rules.md` para enums, paginación y lint del frontend
- Leer `docs/web/pitfalls.md` para errores conocidos (doble request, servicios sin HTTP compartido)

### Instalación de nuevas librerías

Antes de agregar una dependencia (npm/pnpm/Go):

1. **Verificar reputación**: npm trends, GitHub stars, maintenance status, download count
2. **Confirmar que no sea legacy/obsoleta**: buscar alternativas más modernas; si la lib no se actualiza desde hace >1 año, cuestionar su uso
3. **Preferir oficial/maintaineada**: dar prioridad a librerías mantenidas por la organización del framework (ej: `@tanstack/vue-query` sobre alternativas third-party)
4. **Comparar bundle size**: para frontend, usar `bundlephobia.com` o similar
5. **Usar Context7 MCP** para consultar documentación actualizada antes de integrar
6. **Documentar la decisión**: agregar un comentario breve si la elección no es obvia

### Type-check

`vue-tsc --build` es obligatorio tras cualquier cambio en `.ts`/`.vue`. No confiar solo en el linter.

### Swagger

`make swag` tras crear, modificar o eliminar endpoints. Los archivos generados están en `apps/api/docs/`.

## Arquitectura

| Ruta | Contenido |
|---|---|
| `apps/api/` | Backend Go — entrypoint `cmd/server/main.go`, dominios en `internal/domain/<name>/` |
| `apps/web/` | Frontend Vue 3 — entrypoint `src/main.ts` |
| `docs/` | Documentación detallada por tema |
| `.agents/skills/` | Skills instalados (12 disponibles) |

## Convenciones clave

- **Enums**: `const` objects, NO TypeScript `enum`
- **Auth**: Cookies HttpOnly (`access_token` 20 min, `refresh_token` 20 días), no localStorage
- **HTTP client**: `get/post/put/patch/del` desde `@/services/http` (excepción: `healthCatalogService` y `setupService` usan su propio `request()`)
- **Zod**: schemas en `src/schemas/` para formularios (vee-validate + Zod)
- **Contratos FE/BE**: los tipos TypeScript deben reflejar exactamente las structs Go

## Tests

No existen tests Go ni frontend (Vitest no configurado).

## Para empezar a trabajar

1. Leer `docs/architecture.md` si no se conoce la estructura
2. Cargar el skill correspondiente (Vue, Zod, Turborepo, etc.)
3. Usar Context7 MCP para documentación de librerías
4. Tras cambios: `pnpm lint:oxlint` → `pnpm lint:eslint` → `pnpm type-check`
