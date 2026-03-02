# AGENTS.md

Guía para agentes de código (AI) que trabajen en este monorepo.

---

## Estructura del repositorio

```
my_pets_monorepo/
├── apps/
│   ├── api/                  # Backend — Go 1.23+ + Gin
│   │   ├── cmd/server/       # Entrypoint (package main)
│   │   ├── .env              # Variables de entorno locales (no commitear)
│   │   └── internal/
│   │       ├── handlers/     # HTTP handlers (package handlers)
│   │       ├── models/       # Data structs (package models)
│   │       ├── repository/   # Capa de persistencia (interface + postgres)
│   │       └── routes/       # Route registration (package routes)
│   └── web/                  # Frontend — Vue 3 + Vite + TypeScript
│       └── src/
│           ├── types/        # Interfaces TS compartidas
│           ├── services/     # Capa fetch/HTTP
│           ├── stores/       # Pinia stores
│           ├── views/        # Componentes de página
│           ├── components/   # Componentes reutilizables
│           ├── composables/  # Composables Vue
│           └── router/       # Vue Router
├── packages/                 # Paquetes compartidos (futuro)
├── turbo.json
├── package.json              # pnpm workspaces root
├── pnpm-workspace.yaml       # Declaración de workspaces pnpm
└── Makefile                  # Comandos Go
```

---

## Comandos de desarrollo

### Desde la raíz

```bash
pnpm install          # Instalar todas las dependencias JS
pnpm dev              # Dev server frontend (Turborepo)
pnpm build            # Build producción frontend
pnpm lint             # Lint frontend
make docker-dev       # Levantar todo con Docker (API + web, hot-reload)
make docker-down      # Detener Docker
```

### Backend Go

```bash
make dev-api          # Hot-reload con air (.air.toml)
make build-api        # Compilar → bin/server
make test-api         # go test ./... -v
make lint-api         # go vet ./...
make tidy             # go mod tidy
```

#### Correr un único test Go

```bash
# Un test específico por nombre:
cd apps/api && go test ./internal/handlers/ -run TestGetPets -v

# Un paquete específico:
cd apps/api && go test ./internal/handlers/ -v

# Todos los tests:
cd apps/api && go test ./... -v
```

### Frontend (desde apps/web o vía turbo)

```bash
# Desde la raíz:
pnpm dev              # vite dev server en :3000
pnpm build            # type-check + vite build
pnpm lint             # oxlint --fix → eslint --fix

# Desde apps/web:
pnpm type-check       # vue-tsc --build
pnpm lint:oxlint      # oxlint . --fix  (rápido, Rust)
pnpm lint:eslint      # eslint . --fix --cache
pnpm build-only       # vite build (sin type-check)
```

---

## Convenciones — Go

### Estructura de paquetes
- `cmd/server/main.go` → package `main` (solo arranque, config, inyección de dependencias)
- `internal/` → paquetes privados nombrados igual que su carpeta
- No crear paquetes fuera de `cmd/` o `internal/` salvo justificación clara

### Naming
- Funciones exportadas: `PascalCase` (`GetPets`, `CreatePet`)
- Variables/funciones no exportadas: `camelCase` (`nextID`, `pets`)
- Handlers nombrados como `VerbNoun`: `GetPets`, `UpdatePet`, `DeletePet`
- Constantes: `PascalCase` o `SCREAMING_SNAKE` solo si es una constante global conocida

### Imports
Agrupar en bloques separados por línea en blanco:
```go
import (
    "net/http"       // 1. stdlib
    "strconv"

    "github.com/gin-gonic/gin"            // 2. terceros

    "github.com/my-pets/api/internal/models" // 3. internos
)
```

### Structs y tags
```go
type Pet struct {
    ID      uint   `json:"id"`
    Name    string `json:"name" binding:"required"`
    Species string `json:"species" binding:"required"`
    Breed   string `json:"breed"`
}
```
- Siempre incluir `json:` tags en snake_case
- Usar `binding:"required"` en campos obligatorios para validación de Gin
- Timestamps como `time.Time` con tags `json:"created_at"`

### Manejo de errores
- **Early return**: verificar error inmediatamente, responder JSON y `return`
- Errores de cliente → `http.StatusBadRequest` + `gin.H{"error": "..."}`
- Recurso no encontrado → `http.StatusNotFound` + `gin.H{"error": "pet not found"}`
- `log.Fatalf` **solo** en `main.go` para errores de arranque
- No usar `panic()` en handlers

```go
id, err := strconv.Atoi(c.Param("id"))
if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
    return
}
```

### Respuestas de la API
```go
// Lista:   gin.H{"data": slice, "total": len(slice)}
// Recurso: gin.H{"data": pet}
// Delete:  gin.H{"message": "pet deleted"}
// Error:   gin.H{"error": "descripción"}
```

---

## Convenciones — TypeScript / Vue

### Componentes Vue
- **Solo Composition API** con `<script setup lang="ts">` — prohibido Options API
- Orden de bloques: `<script setup>` → `<template>` → `<style scoped>`
- Siempre `<style scoped>` en componentes (no estilos globales en componentes)
- Vistas sufijadas con `View`: `PetsView.vue`, `PetDetailView.vue`
- Archivos en `PascalCase`

### Imports
```ts
// ✅ Correcto
import { ref, onMounted } from 'vue'
import type { Pet, PetPayload } from '@/types/pet'
import { usePetStore } from '@/stores/pets'

// ❌ Evitar imports relativos excepto en router/index.ts
import { usePetStore } from '../stores/pets'
```
- Usar siempre el alias `@/` para imports internos
- `import type { ... }` para importar solo tipos (requerido por el linter)
- Nunca importar desde `vue-router` o `pinia` directamente en componentes si hay un store o composable que lo encapsule

### Tipos e interfaces
```ts
// types/pet.ts — patrón a seguir
export interface Pet { ... }                        // entidad completa (viene del servidor)
export interface PetPayload { ... }                  // body de create/update (sin campos servidor)
export interface ApiResponse<T> { data: T; total?: number }  // envoltorio genérico
```
- Usar `interface` para formas de objeto (no `type` alias)
- Separar siempre el payload de creación/edición del tipo de entidad completo
- Campos opcionales con `?` (no `| undefined` explícito)

### Pinia stores
```ts
// Setup store (no Options store)
export const usePetStore = defineStore('pets', () => {
    const pets = ref<Pet[]>([])
    const loading = ref(false)
    const error = ref<string | null>(null)
    // ...
    return { pets, loading, error, fetchPets, createPet, updatePet, deletePet }
})
```
- Nombre del archivo: lowercase plural (`pets.ts`)
- Nombre del store: `use` + PascalCase + `Store` (`usePetStore`)
- Siempre exponer `loading` y `error` como state reactivo en stores de recursos
- try/catch en operaciones de lectura; dejar propagar errores en mutaciones

### Capa de servicios
```ts
// Patrón: objeto exportado, no clase
export const petService = {
    getAll(): Promise<ApiResponse<Pet[]>> { ... },
    create(payload: PetPayload): Promise<ApiResponse<Pet>> { ... },
}
```
- Archivo: `src/services/[entidad]Service.ts`
- Usar `fetch` nativo (no Axios)
- `BASE_URL = '/api/v1'` — el proxy de Vite redirige a la API en dev
- Helper genérico `request<T>()` para no repetir headers y manejo de errores HTTP

### Naming frontend
| Elemento | Convención | Ejemplo |
|---|---|---|
| Archivos Vue | PascalCase | `PetDetailView.vue` |
| Vistas | sufijo `View` | `PetsView.vue` |
| Stores | `use...Store` | `usePetStore` |
| Servicios | `...Service` | `petService` |
| Handlers en vistas | prefijo `handle` | `handleCreate`, `handleDelete` |
| Rutas nombradas | kebab-case | `'pet-detail'` |

### Router
- Usar `createWebHistory` (HTML5 history mode)
- `HomeView` cargada de forma eager; resto de vistas con **lazy-load**: `() => import(...)`

---

## Contrato API

El frontend espera y la API devuelve:
```
GET  /api/v1/pets       → { data: Pet[], total: number }
GET  /api/v1/pets/:id   → { data: Pet }
POST /api/v1/pets       → { data: Pet }  (body: PetPayload)
PUT  /api/v1/pets/:id   → { data: Pet }  (body: PetPayload)
DELETE /api/v1/pets/:id → { message: string }
GET  /health            → { status: "ok", service: "my-pets-api" }
Errores                 → { error: string }
```

---

## Notas importantes para agentes

1. **Base de datos PostgreSQL** — el backend usa `internal/repository/` con una interface `PetRepository`. La implementación está en `postgres.go`. Credenciales en `apps/api/.env` (no commitear). La migración de la tabla `pets` corre automáticamente al arrancar.
2. **No hay autenticación** — agregar middleware de Gin antes de añadir rutas protegidas.
3. **No hay tests en el frontend** — Vitest no está configurado; no asumir que existe.
4. **MongoDB driver presente** como dependencia indirecta de Gin (no usarlo sin configurarlo explícitamente).
5. **oxlint corre antes que ESLint** (`lint:oxlint` → `lint:eslint`). Ambos con `--fix`. No saltarse ninguno.
6. `noUncheckedIndexedAccess: true` — accesos a arrays/objetos por índice devuelven `T | undefined`; manejar siempre ese caso.
