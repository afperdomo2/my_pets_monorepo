# 🐾 My Pets Monorepo

Monorepo para gestión de mascotas con **Go + Gin** (API REST) y **Vue 3 + TypeScript** (frontend).

## 📁 Estructura

```
my_pets_monorepo/
├── apps/
│   ├── api/        # Backend — Go + Gin
│   └── web/        # Frontend — Vue 3 + Vite + TS
├── packages/       # Código compartido (futuro)
├── docker-compose.yml
├── turbo.json
└── Makefile
```

## ⚙️ Requisitos

| Herramienta | Versión mínima |
|---|---|
| Go | 1.23+ |
| Node.js | 22+ |
| pnpm | 10+ |
| Docker + Compose | Cualquier versión reciente |

---

## 🚀 Desarrollo local

### Opción A — Docker (recomendado)

Levanta API y frontend con hot-reload en un solo comando:

```bash
make docker-dev
```

| Servicio | URL |
|---|---|
| Frontend | <http://localhost:3000> |
| API | <http://localhost:8080> |

### Opción B — Sin Docker

**1. Instalar dependencias JS:**

```bash
pnpm install
```

**2. Backend (requiere [`air`](https://github.com/air-verse/air) para hot-reload):**

```bash
make dev-api
```

**3. Frontend (en otra terminal):**

```bash
pnpm dev
```

> El frontend corre en `:3000` y hace proxy de `/api/*` → `:8080` automáticamente via Vite.

---

## 🔨 Comandos disponibles

### Raíz del monorepo

| Comando | Descripción |
|---|---|
| `pnpm dev` | Dev server del frontend (Turborepo) |
| `pnpm build` | Build de producción del frontend |
| `pnpm lint` | Lint del frontend |
| `make docker-dev` | Levantar todo con Docker (dev) |
| `make docker-down` | Detener Docker |

### Backend Go (`make`)

| Comando | Descripción |
|---|---|
| `make dev-api` | Hot-reload con `air` |
| `make build-api` | Compilar binario en `bin/server` |
| `make test-api` | Ejecutar tests Go |
| `make lint-api` | `go vet` |
| `make tidy` | `go mod tidy` |

---

**Ejemplo de body para crear/actualizar:**

```json
{
  "name": "Firulais",
  "species": "dog",
  "breed": "Labrador",
  "age": 3,
  "owner": "Juan"
}
```

---

## 🏗️ Build de producción

```bash
# Frontend
pnpm build             # genera apps/web/dist/

# Backend
make build-api         # genera bin/server
```
