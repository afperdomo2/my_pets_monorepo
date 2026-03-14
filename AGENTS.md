# AGENTS.md

Guía para agentes de código (AI) que trabajen en este monorepo.

---

## Estructura del repositorio

```
my_pets_monorepo/
├── apps/
│   ├── api/                  # Backend — Go 1.23+ + Gin + GORM
│   │   ├── cmd/server/       # Entrypoint (package main)
│   │   ├── docs/             # Swagger generado por swag
│   │   ├── .env              # Variables de entorno locales (no commitear)
│   │   └── internal/
│   │       ├── config/       # Config struct cargada desde env
│   │       ├── database/     # GORM Connect() + AutoMigrate()
│   │       ├── middleware/   # CORS + JWT middleware
│   │       ├── models/       # Structs GORM compartidos (Pet, User)
│   │       ├── server/       # Wire-up de rutas y arranque HTTP
│   │       ├── validation/   # Traducción de errores de validación a español
│   │       └── domain/       # Módulos de negocio (endpoints)
│   │           ├── auth/     # JWT login/logout/refresh/me + Google OAuth
│   │           ├── health/   # GET /health
│   │           ├── pet/      # Mascotas (handler, payload, repository, gorm_repo, routes)
│   │           ├── setup/    # Flujo de primer usuario (handler + routes)
│   │           └── user/     # Usuarios (handler, payload, repository, gorm_repo, routes)
│   └── web/                  # Frontend — Vue 3 + Vite + TypeScript
│       └── src/
│           ├── types/        # Interfaces TS compartidas
│           ├── services/     # Capa fetch/HTTP (authService, petService, setupService)
│           ├── stores/       # Pinia stores (auth, pets)
│           ├── views/        # Componentes de página (*View.vue)
│           ├── components/   # Componentes reutilizables
│           ├── composables/  # Composables Vue
│           └── router/       # Vue Router + guards
├── packages/                 # Paquetes compartidos (futuro)
├── turbo.json
├── package.json              # pnpm workspaces root
├── pnpm-workspace.yaml
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

# Build directo (requiere GONOSUMDB por restricciones de red en algunos entornos):
GONOSUMDB="*" go build ./...
```

#### Correr un único test Go

```bash
# Un test específico por nombre:
go test ./internal/domain/auth/ -run TestLogin -v

# Un paquete completo:
go test ./internal/domain/pet/ -v

# Todos los tests:
go test ./... -v
```

> Ejecutar siempre desde `apps/api/`, no desde la raíz.

### Frontend (desde apps/web o vía turbo)

```bash
pnpm type-check       # vue-tsc --build  ← correr siempre después de cambios TS
pnpm lint:oxlint      # oxlint . --fix   (rápido, Rust — corre primero)
pnpm lint:eslint      # eslint . --fix --cache
pnpm build-only       # vite build (sin type-check)
pnpm build            # type-check + vite build
```

> No hay tests en el frontend — Vitest no está configurado.

---

## Convenciones — Go

### Idioma de la documentación

- **Toda la documentación en el código debe estar en español**, incluyendo comentarios, funciones de comentarios (docstrings), y comentarios en los archivos.
- Esto incluye comentarios de funciones, explicaciones de lógica, y cualquier texto que forme parte del código fuente.

### Estructura de paquetes

- `internal/domain/` agrupa todos los módulos de negocio (auth, pet, user, setup, health). Cada uno es un paquete autónomo con su propio `handler.go`, `payload.go`, `repository.go`, `gorm_repo.go` y `routes.go`.
- `internal/models/` contiene los structs GORM compartidos (`Pet`, `User`) con tags `gorm:"..."` y `json:"..."`.
- `internal/database/` expone `Connect(*gorm.DB)` y `Migrate(*gorm.DB)` — centraliza la conexión y el AutoMigrate.
- `internal/middleware/` y `internal/config/` son infraestructura transversal — separados de los dominios.
- `internal/server/server.go` → wire-up de rutas y arranque HTTP. Recibe `*gorm.DB`.
- `cmd/server/main.go` → solo arranque, config e inyección de dependencias.
- No crear paquetes fuera de `cmd/` o `internal/` salvo justificación clara.

### Naming

- Funciones exportadas: `PascalCase` (`GetPets`, `CreatePet`)
- Variables/funciones no exportadas: `camelCase` (`nextID`, `userRepo`)
- Handlers: `VerbNoun` (`GetPets`, `UpdatePet`, `DeletePet`)
- Constantes: `PascalCase`; `SCREAMING_SNAKE` solo para constantes globales conocidas

### Imports

Tres bloques separados por línea en blanco:

```go
import (
    "net/http"    // 1. stdlib
    "strconv"

    "github.com/gin-gonic/gin"               // 2. terceros

    "github.com/my-pets/api/internal/models" // 3. internos
)
```

### Structs y tags

```go
type Pet struct {
    ID        string    `json:"id"      gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Name      string    `json:"name"    gorm:"not null"`
    Species   string    `json:"species" gorm:"not null"`
    Breed     string    `json:"breed"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

- Siempre `json:` tags en snake_case.
- `gorm:` tags en los structs de `internal/models/` para schema management.
- `binding:"required"` en los **payload** structs (en cada dominio) para validación de Gin.
- Timestamps como `time.Time` — GORM los gestiona automáticamente.
- Los payloads HTTP (`PetPayload`, `CreateUserPayload`, etc.) viven en cada dominio (`internal/domain/<name>/payload.go`), no en `internal/models/`.

### Manejo de errores

- **Early return**: verificar error inmediatamente, responder JSON y `return`.
- Errores de cliente → `http.StatusBadRequest` + `gin.H{"error": "..."}`.
- No encontrado → `http.StatusNotFound` + `gin.H{"error": "pet not found"}`.
- `log.Fatalf` **solo** en `main.go` para errores de arranque. Nunca `panic()` en handlers.

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

### Idioma de la documentación

- **Toda la documentación en el código debe estar en español**, incluyendo comentarios, JSDoc, y cualquier texto que forme parte del código fuente.

### Componentes Vue

- **Solo Composition API** con `<script setup lang="ts">` — prohibido Options API.
- Orden de bloques: `<script setup>` → `<template>` → `<style scoped>`.
- Siempre `<style scoped>`; sin estilos globales en componentes.
- Vistas sufijadas con `View`: `PetsView.vue`, `PetDetailView.vue`.
- Archivos en `PascalCase`.

### Imports

```ts
// ✅ Correcto
import { ref, computed } from 'vue'
import type { Pet, PetPayload } from '@/types/pet'
import { usePetStore } from '@/stores/pets'

// ❌ Evitar imports relativos (excepto en router/index.ts)
import { usePetStore } from '../stores/pets'
```

- Alias `@/` para todos los imports internos.
- `import type { ... }` para importar solo tipos (requerido por el linter).

### Tipos e interfaces

```ts
export interface Pet { ... }               // entidad completa (viene del servidor)
export interface PetPayload { ... }         // body de create/update (sin campos servidor)
export interface ApiResponse<T> { data: T; total?: number }
```

- Usar `interface` para formas de objeto, no `type` alias.
- Separar el payload de creación/edición del tipo de entidad completo.
- Campos opcionales con `?` (no `| undefined` explícito).
- `noUncheckedIndexedAccess: true` — accesos a arrays por índice devuelven `T | undefined`; manejar siempre ese caso.

### Pinia stores

```ts
export const usePetStore = defineStore('pets', () => {
    const pets = ref<Pet[]>([])
    const loading = ref(false)
    const error = ref<string | null>(null)
    return { pets, loading, error, fetchPets, createPet, updatePet, deletePet }
})
```

- Archivo: lowercase plural (`pets.ts`, `auth.ts`).
- Nombre: `use` + PascalCase + `Store`.
- Siempre exponer `loading` y `error` como state reactivo.
- `try/catch` en lecturas; propagar errores en mutaciones.

### Capa de servicios

```ts
export const petService = {
    getAll(): Promise<ApiResponse<Pet[]>> { ... },
    create(payload: PetPayload): Promise<ApiResponse<Pet>> { ... },
}
```

- Archivo: `src/services/[entidad]Service.ts`.
- `fetch` nativo (no Axios). Siempre `credentials: 'include'` (cookies HttpOnly).
- `BASE_URL = '/api/v1'` — el proxy de Vite redirige a `:8080` en dev.
- Helper privado `request<T>()` para centralizar headers y manejo de errores HTTP.

### Naming frontend

| Elemento          | Convención      | Ejemplo                      |
|-------------------|-----------------|------------------------------|
| Archivos Vue      | PascalCase      | `PetDetailView.vue`          |
| Vistas            | sufijo `View`   | `PetsView.vue`               |
| Stores            | `use...Store`   | `usePetStore`                |
| Servicios         | `...Service`    | `petService`                 |
| Handlers en vista | prefijo `handle`| `handleCreate`, `handleDelete`|
| Rutas nombradas   | kebab-case      | `'pet-detail'`               |

### Router

- `createWebHistory` (HTML5 history mode).
- `HomeView` eager; resto de vistas con lazy-load: `() => import(...)`.
- `meta: { requiresAuth: true }` en rutas protegidas.
- El guard `beforeEach` cachea el setup status en variables de módulo; usar `resetSetupCache()` (exportada desde `router/index.ts`) si se necesita invalidar tras crear el primer usuario.

---

## Autenticación y autorización

- **Tokens**: HttpOnly cookies (`access_token` 20 min, `refresh_token` 20 días).
- **Rutas públicas**: `/health`, `/swagger/*`, `/api/v1/auth/login`, `/api/v1/auth/refresh`, `/api/v1/auth/google`, `/api/v1/auth/google/callback`, `/api/v1/setup/status`, `/api/v1/setup`.
- **Rutas protegidas**: todo lo demás bajo `/api/v1` requiere el middleware `middleware.JWT(cfg)`.
- `POST /api/v1/users` — solo `is_system_user` puede crear usuarios.
- `PUT/DELETE /api/v1/users/:id` — solo el propio usuario o `is_system_user`.
- El primer usuario creado por cualquier método recibe `is_system_user = true` automáticamente.

## Flujo de setup (primer uso)

- `GET /api/v1/setup/status` → `{ "needs_setup": bool }` — consultado por el router guard.
- `POST /api/v1/setup` → crea el primer usuario administrador; devuelve 409 si ya existen usuarios.
- Frontend: si `needs_setup: true`, toda navegación redirige a `/setup`; tras crear el usuario, llamar a `resetSetupCache()` antes de `router.push`.

---

## Contrato API completo

```
GET    /health                       → { status: "ok", service: "my-pets-api" }
GET    /api/v1/setup/status          → { needs_setup: bool }
POST   /api/v1/setup                 → { data: User }
POST   /api/v1/auth/login            → { data: User, expires_in: number }  + Set-Cookie
POST   /api/v1/auth/logout           → { message: string }                 + clear cookies
POST   /api/v1/auth/refresh          → { data: User, expires_in: number }  + Set-Cookie
GET    /api/v1/auth/me               → { data: User }                       (protegido)
GET    /api/v1/auth/google           → redirect OAuth Google
GET    /api/v1/auth/google/callback  → redirect frontend
GET    /api/v1/pets                  → { data: Pet[], total: number }
GET    /api/v1/pets/:id              → { data: Pet }
POST   /api/v1/pets                  → { data: Pet }
PUT    /api/v1/pets/:id              → { data: Pet }
DELETE /api/v1/pets/:id              → { message: string }
GET    /api/v1/users                 → { data: User[], total: number }
GET    /api/v1/users/:id             → { data: User }
POST   /api/v1/users                 → { data: User }   (solo system_user)
PUT    /api/v1/users/:id             → { data: User }
DELETE /api/v1/users/:id             → { message: string }
Errores                              → { error: string }
```

---

## TanStack Vue Query — comportamiento interno crítico

El frontend usa **TanStack Vue Query** (`@tanstack/vue-query`) para data fetching. Hay comportamientos internos no obvios que hay que conocer para evitar bugs difíciles de diagnosticar:

### Cómo funciona `useBaseQuery` internamente

- `watch(defaultedOptions, updater)` es el watcher interno que sincroniza cambios de `queryKey` con el observer. Cada vez que el `queryKey` reactivo cambia, Vue llama a `updater()` → `observer.setOptions(newOptions)`.
- `refetch()` de la instancia llama `updater()` + `observer.refetch()` — adopta el `queryKey` actual antes de hacer fetch.
- `shouldFetchOptionally` controla si el observer hace un fetch automático al cambiar de key: devuelve `true` si la nueva query está **stale**. Si la query ya está en caché y no está stale, no hay fetch reactivo automático.

### Problema: doble request al refrescar desde página > 1

Si se invalida el caché **antes** de cambiar `page.value`, la query de página 1 queda stale → al cambiar `page.value = 1` el observer detecta el nuevo key stale → `shouldFetchOptionally` devuelve `true` → fetch reactivo automático → luego `refetch()` hace otro → **2 requests**.

### Solución: patrón de 3 pasos para `refresh()` con paginación

```ts
async function refresh() {
  if (refreshing) return
  refreshing = true          // flag JS síncrono — bloquea clics adicionales inmediatamente
  try {
    page.value = 1           // 1. cambiar página SIN invalidar — la query de pág 1 no está stale,
                             //    shouldFetchOptionally → false, cero fetch reactivo
    await nextTick()         // 2. esperar a que Vue flush el watcher interno (watch(defaultedOptions, updater))
                             //    El observer ya apunta a página 1 sin haber hecho fetch
    await queryClient.invalidateQueries({
      queryKey: ['users'],
      refetchType: 'active', // 3. marca stale TODO el caché de ['users'] y refetchea solo el
    })                       //    observer activo (página 1) = exactamente UN request
  } finally {
    refreshing = false
  }
}
```

### Semántica de `refetchType` en `invalidateQueries`

| `refetchType`  | Efecto                                                                 |
|----------------|------------------------------------------------------------------------|
| `'active'`     | Refetchea solo observers activos en ese momento (default)             |
| `'inactive'`   | Refetchea solo observers inactivos                                    |
| `'all'`        | Refetchea todos los observers                                         |
| `'none'`       | Marca stale pero NO dispara ningún refetch automático                 |

> `refetchType: 'none'` **no** suprime el fetch que dispara `shouldFetchOptionally` cuando el `queryKey` cambia. Solo suprime el refetch que `invalidateQueries` haría por sí mismo.

### Patrón load-more acumulativo (mascotas)

Para "Ver más" que acumula resultados en lugar de reemplazarlos, usar un array local separado del resultado de la query:

```ts
const allPets = ref<Pet[]>([])   // acumulador local
const page = ref(1)

watch(queryData, (newData) => {
  if (page.value === 1) allPets.value = newData  // reset en refresh
  else allPets.value = [...allPets.value, ...newData]  // acumular
})
```

El flag síncrono para bloquear el botón de refresh mientras está en curso:

```ts
let refreshing = false   // JS puro, NO ref — síncrono, no necesita reactividad
```

---

## Paginación — convenciones

### Respuesta paginada del backend

```json
{ "data": [...], "total": N, "page": 1, "per_page": 10, "total_pages": 3 }
```

Los handlers de `GET /api/v1/pets` y `GET /api/v1/users` aceptan query params `page` y `per_page`. Valor por defecto: `page=1`, `per_page=10`.

### Tipo `PaginatedResponse<T>` en frontend

```ts
export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  per_page: number
  total_pages: number
}
```

Definido en `apps/web/src/types/pet.ts`.

### Persistencia de `perPage` en usuarios

El store `apps/web/src/stores/ui.ts` expone `usersPerPage` (valores: 10, 25, 50) persistido en `localStorage` directamente, sin plugins externos.

### Componentes de paginación

- `components/ui/AppPagination.vue` — paginador clásico con ellipsis (`...`), máximo 5 páginas visibles, botones primera/anterior/siguiente/última.
- `components/ui/PerPageSelector.vue` — selector "Mostrar X por página".

---

## Notas importantes para agentes

1. **PostgreSQL** — credenciales en `apps/api/.env` (no commitear). Las migraciones de `pets` y `users` corren automáticamente al arrancar via `database.Migrate()` (GORM AutoMigrate).
2. **GONOSUMDB** — algunos entornos con restricciones de red requieren `GONOSUMDB="*"` para `go build`. Los errores LSP de módulos como `golang-jwt/jwt/v5` u `oauth2` suelen ser falsos positivos; verificar con `go build ./...`.
3. **No hay tests en el frontend** — Vitest no está configurado; no asumir que existe.
4. **oxlint corre antes que ESLint** — no saltarse ninguno de los dos pasos de lint.
5. **`pnpm type-check` es obligatorio** — correrlo siempre tras cambios en TypeScript/Vue.
6. **Cookies HttpOnly** — el frontend usa `credentials: 'include'` en todos los fetch; nunca guardar tokens en localStorage.
7. **TanStack Query + paginación** — al implementar refresh con cambio de página, seguir el patrón de 3 pasos documentado arriba para evitar doble request. Ver sección "TanStack Vue Query — comportamiento interno crítico".
8. **NO crear commits automáticamente** — solo crear commits cuando el usuario lo solicite explícitamente. Completar tareas sin hacer commits; el usuario decidirá cuándo commitear.
