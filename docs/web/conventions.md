# Convenciones — Frontend Vue/TS

## Componentes

- Composition API con `<script setup lang="ts">`
- `<style scoped>` siempre
- Vistas con sufijo `View`: `PetsView.vue`, `LoginView.vue`

## Imports

- Alias `@/` para imports internos
- `import type { ... }` para tipos

## Enums

Usar `const` objects, NO TypeScript `enum`:

```ts
export const HealthRecordStatus = { Pending: 'pending' } as const
export type HealthRecordStatusType = (typeof HealthRecordStatus)[keyof typeof HealthRecordStatus]
```

## HTTP client

Usar `get/post/put/patch/del` desde `@/services/http`:

```ts
import { get, post } from '@/services/http'

export const petService = {
  getAll(page = 1, perPage = 10) { return get(`/pets?page=${page}&per_page=${perPage}`) },
}
```

- `BASE_URL = '/api/v1'` (centralizado en `http.ts`)
- `credentials: 'include'` para cookies JWT

**Excepción**: `healthCatalogService` y `setupService` usan su propio `request()` interno, no el compartido.

## Zod schemas

Validación de formularios en `src/schemas/` (vee-validate + Zod).
Schemas existentes: `pet.ts`, `user.ts`, `healthCatalog.ts`.

## TanStack Vue Query

```ts
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 60 * 1000,       // 60s
      gcTime: 5 * 60 * 1000,      // 5min
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})
```

## Pinia stores

Exponer `loading` y `error` en todas las stores.

## Sesión

`authStore.initSession()` se llama **antes** de `app.mount()` en `main.ts` para restaurar la sesión desde cookies y evitar race conditions en el router guard.
