# Arquitectura del monorepo

## Estructura general

```
my_pets_monorepo/
├── apps/
│   ├── api/        # Backend — Go 1.25 + Gin + GORM
│   └── web/        # Frontend — Vue 3 + Vite + TS + TanStack Query
├── docs/           # Documentación para agentes
├── docker-compose.yml
├── Makefile        # Comandos Go y Docker
├── package.json    # pnpm workspaces root + Turborepo
├── pnpm-workspace.yaml
└── turbo.json
```

## Backend API (`apps/api/`)

**Entrypoint:** `cmd/server/main.go`

**Dominios** (en `internal/domain/<name>/`):
- `auth` — login local, Google OAuth, JWT, refresh
- `dashboard` — resumen de métricas
- `exam` — exámenes veterinarios CRUD
- `health` — liveness check (`GET /health`)
- `health_catalog` — catálogo de vacunas/desparasitaciones
- `health_record` — registros de salud de mascotas
- `pet` — mascotas CRUD
- `setup` — configuración inicial (primer usuario)
- `user` — usuarios CRUD
- `vaccine_application` — aplicaciones de vacunas
- `profile/` — **vacío** (no implementado)

**Paquetes internos compartidos:** `config/`, `database/`, `middleware/`, `models/`, `server/`, `validation/`

**Cableado:** `internal/server/server.go` crea repos compartidos (`userRepo`), instancia handlers por dominio, y registra rutas en grupos público (`/api/v1`, sin JWT) y protegido (`/api/v1` con middleware JWT).

## Frontend web (`apps/web/`)

**Entrypoint:** `src/main.ts`

**Estructura de `src/`:**
- `composables/` — hooks TanStack Query (CRUD + cache invalidation)
- `services/` — llamadas HTTP
- `stores/` — Pinia (sesión, UI)
- `types/` — interfaces TypeScript
- `constants/` — constantes con `const` objects (no TS enum)
- `utils/` — funciones puras
- `schemas/` — validación Zod para formularios (vee-validate + Zod)
- `views/` — páginas, sufijo `View`
- `components/` — componentes reutilizables

## Comandos del monorepo (raíz)

| Comando | Descripción |
|---|---|
| `pnpm install` | Instalar dependencias JS |
| `pnpm dev` | Dev server frontend (Turborepo) |
| `pnpm build` | Build producción frontend (Turborepo) |
| `pnpm lint` | Lint frontend (Turborepo) |
| `make docker-dev` | Levantar Postgres + API + frontend |
| `make docker-down` | Detener Docker |

## Autenticación

- Cookies HttpOnly: `access_token` (20 min), `refresh_token` (20 días)
- No se usa localStorage
- `authStore.initSession()` se llama **antes** de `app.mount()` en `main.ts`

## Docker Compose

Tres servicios: `db` (Postgres 16-alpine), `api` (air hot-reload), `web` (vite dev).
Frontend en `:3000` con proxy de `/api/*` → API `:8080`.
