# Comandos — Frontend Vue/TS

Ejecutar desde `apps/web/`.

## Desarrollo

| Comando | Descripción |
|---|---|
| `pnpm dev` | Vite dev server en `:3000`, proxy `/api/*` → `:8080` |
| `pnpm build` | `run-p type-check "build-only"` |
| `pnpm type-check` | `vue-tsc --build` (**obligatorio** tras cambios `.ts`/`.vue`) |
| `pnpm lint:oxlint` | `oxlint . --fix` (ejecutar **primero**) |
| `pnpm lint:eslint` | `eslint . --fix --cache` (ejecutar **después**) |
| `pnpm preview` | Vista previa del build |

## Tests

> No hay tests frontend (Vitest no configurado).
