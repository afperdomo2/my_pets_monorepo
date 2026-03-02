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
| Node.js | 20+ |
| npm | 10+ |
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
npm install
```

**2. Backend (requiere [`air`](https://github.com/air-verse/air) para hot-reload):**

```bash
make dev-api
```

**3. Frontend (en otra terminal):**

```bash
npm run dev
```

> El frontend corre en `:3000` y hace proxy de `/api/*` → `:8080` automáticamente via Vite.

---

## 🔨 Comandos disponibles

### Raíz del monorepo

| Comando | Descripción |
|---|---|
| `npm run dev` | Dev server del frontend (Turborepo) |
| `npm run build` | Build de producción del frontend |
| `npm run lint` | Lint del frontend |
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

## 🌐 API Endpoints

Base URL: `http://localhost:8080/api/v1`

| Método | Ruta | Descripción |
|---|---|---|
| `GET` | `/pets` | Listar mascotas |
| `GET` | `/pets/:id` | Obtener mascota |
| `POST` | `/pets` | Crear mascota |
| `PUT` | `/pets/:id` | Actualizar mascota |
| `DELETE` | `/pets/:id` | Eliminar mascota |
| `GET` | `/health` | Health check |

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
npm run build          # genera apps/web/dist/

# Backend
make build-api         # genera bin/server
```
