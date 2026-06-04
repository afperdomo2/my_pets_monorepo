# Reglas del proyecto

## Generales

- **Documentación en español** — comentarios, docstrings, commits, PRs
- **NO commits automáticos** — solo cuando el usuario lo pida explícitamente
- **Contratos frontend-backend**: los tipos TypeScript deben reflejar exactamente las structs Go

## Frontend

- **Enums**: usar `const` objects, NO TypeScript `enum`
  ```ts
  export const Status = { Active: 'active' } as const
  export type StatusType = (typeof Status)[keyof typeof Status]
  ```
- **Paginación**: usar `PER_PAGE_DEFAULT = 10` desde `@/constants/pagination`
- **Lint**: `pnpm lint:oxlint` **primero**, `pnpm lint:eslint` después
- **Type-check**: obligatorio tras cualquier cambio en `.ts`/`.vue`

## Backend

- **Swagger**: `make swag` tras cambios en handlers o rutas
- **Migraciones**: GORM AutoMigrate corre al iniciar (`internal/database/`)
- **Redes restringidas**: usar `GONOSUMDB="*" go build ./...` si hay problemas de red
