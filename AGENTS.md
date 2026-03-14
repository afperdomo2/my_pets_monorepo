# AGENTS.md

Guía para agentes de código (AI) que operen en este monorepo.

---

## Estructura del repositorio

```
my_pets_monorepo/
├── apps/
│   ├── api/                  # Backend — Go 1.23+ + Gin + GORM
│   │   ├── cmd/server/       # Entrypoint (main.go)
│   │   ├── docs/             # Swagger generado por swag
│   │   └── internal/
│   │       ├── config/, database/, middleware/, models/, server/, validation/
│   │       └── domain/       # Módulos de negocio (auth, pet, user, health_catalog, health_record, setup)
│   └── web/                  # Frontend — Vue 3 + Vite + TypeScript + TanStack Query
│       └── src/
│           ├── types/, services/, stores/, views/, components/, composables/, router/
├── package.json              # pnpm workspaces root
├── pnpm-workspace.yaml
└── Makefile                  # Comandos Go
```

---

## Comandos de desarrollo

### Root

```bash
pnpm install      # Instalar dependencias JS
pnpm dev          # Dev server frontend
pnpm build        # Build producción
pnpm lint         # Lint frontend
make docker-dev   # Levantar todo con Docker
make docker-down  # Detener Docker
```

### Backend Go

```bash
cd apps/api/

make dev-api       # Hot-reload con air
make build-api     # Compilar → bin/server
make test-api      # go test ./... -v
make lint-api      # go vet ./...
make tidy          # go mod tidy

# Un test específico:
go test ./internal/domain/pet/ -run TestCreatePet -v
go test ./... -v   # Todos los tests
```

**Build con restricciones de red:** `GONOSUMDB="*" go build ./...`

### Frontend

```bash
cd apps/web/

pnpm type-check    # vue-tsc --build (OBLIGATORIO tras cambios TS)
pnpm lint:oxlint   # oxlint . --fix (primero)
pnpm lint:eslint   # eslint . --fix --cache
pnpm build         # type-check + vite build
```

> No hay tests en frontend (Vitest no configurado).

---

## Convenciones — Go

### Documentación
- **Toda la documentación en español**: comentarios, docstrings, etc.

### Estructura de paquetes
- `internal/domain/<name>/` → paquete autónomo: `handler.go`, `payload.go`, `repository.go`, `gorm_repo.go`, `routes.go`
- `internal/models/` → structs GORM compartidos
- `cmd/server/main.go` → solo arranque

### Imports (3 bloques)
```go
import (
    "net/http"                    // 1. stdlib
    "github.com/gin-gonic/gin"    // 2. terceros
    "github.com/my-pets/api/internal/models" // 3. internos
)
```

### Structs y tags
```go
type Pet struct {
    ID        string    `json:"id"      gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Name      string    `json:"name"    gorm:"not null"`
}
```
- `json:` siempre snake_case
- `binding:"required"` en payload structs
- Timestamps como `time.Time`

### Naming
- Exportadas: `PascalCase` | No exportadas: `camelCase`
- Handlers: `VerbNoun` (`GetPets`, `UpdatePet`)

### Manejo de errores
- Early return: verificar error, responder JSON y return
- Errores cliente → `http.StatusBadRequest` + `gin.H{"error": "..."}`
- No encontrado → `http.StatusNotFound`
- `log.Fatalf` solo en main.go, nunca `panic()` en handlers

### Respuestas API
```go
// Lista:   gin.H{"data": slice, "total": N, "page": 1, "per_page": 10, "total_pages": 3}
// Recurso: gin.H{"data": pet}
// Delete:  gin.H{"message": "registro eliminado"}
// Error:   gin.H{"error": "descripción"}
```

---

## Convenciones — TypeScript / Vue

### Documentación
- **Toda la documentación en español**.

### Componentes Vue
- Solo Composition API con `<script setup lang="ts">`
- Orden: `<script setup>` → `<template>` → `<style scoped>`
- Siempre `<style scoped>`
- Vistas con sufijo `View`: `PetsView.vue`

### Imports
- Alias `@/` para imports internos
- `import type { ... }` para tipos

### Tipos y Contratos

**Estructura de tipos:**
```
src/types/
├── shared.ts          # ApiResponse<T>, PaginatedResponse<T> (reutilizables)
├── pet.ts             # Pet, CreatePetPayload, UpdatePetPayload
├── user.ts            # User, UserWithPetCount, payloads
├── healthRecord.ts    # HealthRecord, payloads
├── healthCatalog.ts   # HealthCatalog, payloads
├── auth.ts            # LoginCredentials, LoginResponse, MeResponse
└── index.ts           # Índice de exportación
```

**Interfaces estándar:**
```ts
export interface Pet { }              // Entidad completa
export interface PetPayload { }       // Body create/update
export interface ApiResponse<T> { data: T; total?: number }
export interface PaginatedResponse<T> {
  data: T[]; total: number; page: number; per_page: number; total_pages: number;
}
```

**Enums (patrón objeto const):**
- Usar `const` objects en lugar de `enum` de TypeScript para mejor type safety y strings puros en runtime
- Los enums deben reflejar exactamente los valores del backend Go

```ts
// ✅ Patrón recomendado (Opción B)
export const HealthRecordStatus = {
  Pending: 'pending',
  Applied: 'applied',
  Overdue: 'overdue',
} as const

export type HealthRecordStatusType = (typeof HealthRecordStatus)[keyof typeof HealthRecordStatus]

// Uso: HealthRecordStatus.Pending → 'pending'
```

**Constantes centralizadas:**
```
src/constants/
├── pagination.ts      # PER_PAGE_DEFAULT = 10, PER_PAGE_OPTIONS, helpers
├── petSize.ts         # PET_SIZE_VALUES, labels, icons
├── species.ts         # PET_SPECIES, helpers
├── lifeStage.ts       # LifeStageForDog, LifeStageForCat, etc. + labels
├── healthRecord.ts    # HealthRecordStatus, HealthCatalogCategory + labels
├── auth.ts            # AuthProvider, token expiry, rutas
└── index.ts           # Índice de exportación
```

**Utilidades y helpers:**
```
src/utils/
├── date.ts            # formatDate, isBirthday, calculateAge, estimateBirthDate
├── formatters.ts      # formatWeight, formatAge, capitalize, truncate
├── healthRecord.ts    # getHealthRecordStatusLabel, isOverdue, daysUntilDue
├── pet.ts             # Funciones específicas de pets (compatibilidad)
└── index.ts           # Índice de exportación
```

### Pinia stores
```ts
export const usePetStore = defineStore('pets', () => {
    const pets = ref<Pet[]>([])
    const loading = ref(false)
    const error = ref<string | null>(null)
    return { pets, loading, error, fetchPets, createPet }
})
```
- Exponer siempre `loading` y `error`

### Servicios

**Cliente HTTP compartido (`src/services/http.ts`):**
```ts
import { get, post, put, patch, del } from '@/services/http'

export const petService = {
  getAll(page = 1, perPage = PER_PAGE_DEFAULT) {
    return get(`/pets?page=${page}&per_page=${perPage}`)
  },
  // ...
}
```

- `BASE_URL = '/api/v1'` (centralizada en http.ts)
- `credentials: 'include'` para cookies JWT
- `ApiError` class para manejo uniforme de errores

**Estructura de servicios:**
```
src/services/
├── http.ts              # Cliente HTTP compartido (get, post, put, patch, del)
├── authService.ts       # login, logout, refresh, me, updateProfile
├── petService.ts        # getAll, getById, create, update, remove
├── userService.ts       # getAll, getById, create, update, remove
├── healthRecordService.ts
└── healthCatalogService.ts
```

---

## Rutas API

### Públicas
```
GET  /health
GET  /api/v1/setup/status
POST /api/v1/setup
POST /api/v1/auth/login
POST /api/v1/auth/refresh
GET  /api/v1/auth/google
GET  /api/v1/auth/google/callback
```

### Protegidas (requieren JWT cookie)
```
GET    /api/v1/auth/me
GET    /api/v1/pets
GET    /api/v1/pets/:id
POST   /api/v1/pets
PUT    /api/v1/pets/:id
DELETE /api/v1/pets/:id

GET    /api/v1/health-records
POST   /api/v1/health-records
PUT    /api/v1/health-records/:record_id
PATCH  /api/v1/health-records/:record_id/status
DELETE /api/v1/health-records/:record_id
GET    /api/v1/health-records/pets/:pet_id
GET    /api/v1/health-records/pets/:pet_id/category/:category

GET    /api/v1/health-catalogs/category/:category
GET    /api/v1/health-catalogs/:id
POST   /api/v1/health-catalogs (solo system_user)
PUT    /api/v1/health-catalogs/:id (solo system_user)
DELETE /api/v1/health-catalogs/:id (solo system_user)

GET    /api/v1/users
GET    /api/v1/users/:id
POST   /api/v1/users (solo system_user)
PUT    /api/v1/users/:id
DELETE /api/v1/users/:id
```

### Contraseña y tokens
- Cookies HttpOnly: `access_token` (20 min), `refresh_token` (20 días)
- No usar localStorage

---

## TanStack Vue Query — Errores comunes

### Doble request al hacer refresh con paginación

**Problema:** Invalidar caché antes de cambiar página → stale query → fetch reactivo + refetch() = 2 requests.

**Solución:** Patrón de 3 pasos:
```ts
async function refresh() {
  if (refreshing) return
  refreshing = true
  try {
    page.value = 1           // 1. Cambiar página SIN invalidar
    await nextTick()         // 2. Esperar watcher interno
    await queryClient.invalidateQueries({
      queryKey: ['entidad'],
      refetchType: 'active', // 3. Refetch solo observer activo
    })
  } finally {
    refreshing = false
  }
}
```

---

## Notas importantes

1. **GONOSUMDB**: Algunos entornos requieren `GONOSUMDB="*"` para build
2. **Migraciones**: GORM AutoMigrate corre al iniciar (definidas en `internal/database/`)
3. **Swagger**: Regenerar con `make swag` tras cambios en handlers/routes
4. **Type-check**: Obligatorio tras cambios TS (`pnpm type-check`)
5. **Lint**: `oxlint` antes de `eslint`
6. **NO commits automáticos**: Solo cuando el usuario lo pida explícitamente
7. **Contratos frontend-backend**: Los tipos TypeScript deben reflejar exactamente las structs Go del backend
8. **Enums**: Usar patrón `const` object (Opción B) en lugar de `enum` TS — ver sección "Convenciones — TypeScript / Vue"
9. **Paginación**: Usar `PER_PAGE_DEFAULT = 10` desde `@/constants/pagination` en todos los composables
10. **HTTP client**: Usar `get/post/put/patch/del` desde `@/services/http` en lugar de `fetch` directo
