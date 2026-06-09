# 🐾 My Pets Monorepo

> Gestión de mascotas — API REST con Go + Gin y frontend Vue 3 + TypeScript.

---

## 📁 Estructura

```
my_pets_monorepo/
├── apps/
│   ├── api/              # Backend — Go 1.25 + Gin + GORM
│   └── web/              # Frontend — Vue 3 + Vite 7 + TS
├── .github/workflows/    # CI — GitHub Actions
├── docs/
│   └── api/
│       ├── commands.md   # Comandos de backend
│       ├── conventions.md# Convenciones de código
│       └── testing.md    # Guía de testing (227 tests)
├── docker-compose.prod.yml  # Producción (build estático)
├── lefthook.yml             # Pre-commit hooks
├── Makefile                 # Comandos del backend
├── .env.example             # Variables de entorno
└── turbo.json               # Turborepo config
```

---

## ⚙️ Requisitos

| Herramienta | Versión mínima |
|---|---|
| Go | 1.25+ |
| Node.js | 22+ |
| pnpm | 10+ |
| Docker + Compose | Cualquier versión reciente (solo para producción) |

---

## 🚀 Desarrollo local

```bash
# 1. Instalar dependencias JS
pnpm install

# 2. Backend (requiere air para hot-reload)
make dev-api

# 3. Frontend (otra terminal)
pnpm dev
```

> El frontend corre en `:3000` y hace proxy de `/api/*` → `:8080` automáticamente via Vite.

---

## 🚢 Producción con Docker

### 🖥️ Local (PostgreSQL en contenedor)

```bash
# 1. Configurar variables de entorno
cp .env.example .env
# Editar JWT_SECRET, POSTGRES_PASSWORD y DATABASE_URL

# 2. Build y despliegue
docker compose -f docker-compose.prod.yml up --build -d

# 3. Detener servicios
docker compose -f docker-compose.prod.yml down
```

---

## 🔨 Comandos disponibles

### 🖥️ Backend (desde raíz, vía Makefile)

| Comando | Descripción | Docker |
|---|---|---|
| `make dev-api` | Hot-reload con air (local) | ❌ |
| `make build-api` | Compilar binario `bin/server` | ❌ |
| `make test-api` | Tests handler + unit (~0.1s) | ❌ |
| `make test-api-integration` | Tests handler + unit + repo (~5 min) | ✅ |
| `make lint-api` | `go vet ./...` | ❌ |
| `make tidy` | `go mod tidy` | ❌ |
| `make swag` | Regenerar Swagger docs | ❌ |
| `make help` | Listar todos los targets | ❌ |

Ver detalle en [`docs/api/commands.md`](docs/api/commands.md).

### 🌐 Frontend (desde `apps/web/`)

| Comando | Descripción |
|---|---|
| `pnpm dev` | Vite dev server `:3000` |
| `pnpm build` | Build producción |
| `pnpm lint` | Oxlint + ESLint |
| `pnpm type-check` | `vue-tsc --build` |

### 📦 Raíz del monorepo

| Comando | Descripción |
|---|---|
| `pnpm dev` | Turborepo dev (frontend) |
| `pnpm build` | Turborepo build |
| `pnpm lint` | Turborepo lint |

---

## 🧪 Testing

**227 tests** — 3 tipos:

| Tipo | Stack | Cantidad |
|---|---|---|
| Handler tests | `testify/mock` + `gin.CreateTestContext` | 130 |
| Pure unit tests | `testing` + `testify/require` | 43 |
| Repository tests | `testcontainers-go` + Postgres 16 | 54 |

Guía completa en [`docs/api/testing.md`](docs/api/testing.md).

---

## 🤖 CI — GitHub Actions

| Workflow | Trigger | Se ejecuta si cambia |
|---|---|---|
| [`backend.yml`](.github/workflows/backend.yml) | Push/PR a `main` o `develop` | `apps/api/**` |
| [`frontend.yml`](.github/workflows/frontend.yml) | Push/PR a `main` o `develop` | `apps/web/**` |

- Backend: `go vet` + `go test` (handler + unit, sin Docker)
- Frontend: `pnpm install` → `lint` → `type-check` → `build`
- Los repo tests con Docker solo se ejecutan localmente bajo pedido

---

## 🐶 Pre-commit hooks (Lefthook)

Al hacer `git commit`, se ejecutan validaciones según los archivos modificados:

| Job | Se ejecuta si cambian |
|---|---|
| `lint-api` | `apps/api/**/*.go` |
| `test-api` | `apps/api/**/*.go` |
| `test-api-integration` | `apps/api/**/gorm_repo.go` |
| `lint-web` | `apps/web/**/*.{ts,vue}` |
| `type-check` | `apps/web/**/*.{ts,vue}` |

Config: [`lefthook.yml`](lefthook.yml)

---

## 📚 Documentación detallada

| Recurso | Descripción |
|---|---|
| [`docs/api/commands.md`](docs/api/commands.md) | Comandos de desarrollo y testing |
| [`docs/api/conventions.md`](docs/api/conventions.md) | Convenciones de código Go |
| [`docs/api/testing.md`](docs/api/testing.md) | Guía completa de testing (227 tests) |
| [`AGENTS.md`](AGENTS.md) | Guía para agentes IA |

---

## 📖 Endpoints — Swagger

```
http://localhost:8080/swagger/index.html
```

Regenerar tras cambios:

```bash
make swag
```

---

## 🏗️ Build manual (sin Docker)

```bash
# Frontend
pnpm build           # genera apps/web/dist/

# Backend
make build-api       # genera bin/server
```
