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
├── docker-compose.prod.yml  # Producción local (build estático)
├── docker-compose.aws.yml   # AWS EC2 (imágenes desde ECR)
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

# 2. Editar las variables obligatorias en .env
```

| Variable | Descripción | Ejemplo |
|---|---|---|
| `POSTGRES_PASSWORD` | Contraseña de PostgreSQL | `change-this-password-in-production` |
| `DATABASE_URL` | DSN de conexión a la DB | `host=db user=my_pets_user password=... dbname=my_pets port=5432 sslmode=disable` |
| `JWT_SECRET` | Secreto para firmar tokens JWT | `change-this-to-a-random-secret-in-production` |

> `GOOGLE_CLIENT_ID` y `GOOGLE_CLIENT_SECRET` son opcionales (solo se necesitan para OAuth).

```bash
# 3. Build y despliegue
docker compose -f docker-compose.prod.yml up --build -d

# 3. Build y despliegue (separado)
docker compose -f docker-compose.prod.yml build

docker compose -f docker-compose.prod.yml up -d

# 4. Detener servicios
docker compose -f docker-compose.prod.yml down
```

---

## ☁️ AWS EC2 con ECR

Deploy en un EC2 usando imágenes pre-construidas desde Amazon ECR.
**No se compila en el servidor** — se evita el OOM por RAM limitada (t3.micro).

### 📋 Prerrequisitos

| Recurso | Descripción |
|---|---|
| Repos ECR | `my-pets-api` y `my-pets-web` creados en AWS |
| GitHub Secrets | `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION` |
| EC2 | Docker instalado y autenticado contra ECR |

**Ruta:** Settings → Secrets and variables → Actions → Repository secrets (la pestaña por defecto, no "Environment secrets").

```bash
# Crear repos en ECR (una vez)
aws ecr create-repository --repository-name my-pets-api
aws ecr create-repository --repository-name my-pets-web
```

### 🏗️ Build & push manual

Desde tu máquina local (no en EC2):

```bash
export AWS_ACCOUNT_ID=123456789012
export AWS_REGION=us-east-1
export IMAGE_TAG=$(git rev-parse --short HEAD)

bash scripts/build-and-push.sh
```

Esto construye las imágenes, las tagea con el SHA del commit + `latest`, y las sube a ECR.

### 🤖 CI/CD automático

El workflow [`deploy-ecr.yml`](.github/workflows/deploy-ecr.yml) se activa al pushear a `main` cuando hay cambios en `apps/api/`, `apps/web/`, `pnpm-lock.yaml` o `pnpm-workspace.yaml`.

Build + push automáticos sin intervención manual.

### 🚀 Deploy en EC2

```bash
# 1. Autenticar contra ECR (una vez por sesión)
aws ecr get-login-password --region us-east-1 \
  | docker login --username AWS --password-stdin $ACCOUNT.dkr.ecr.$REGION.amazonaws.com

# 2. Configurar variables
cp .env.example .env
# Editar POSTGRES_PASSWORD, JWT_SECRET, APP_URL, FRONTEND_URL
# APP_URL=http://<ip-publica-ec2>
# FRONTEND_URL=http://<ip-publica-ec2>

# 3. Descargar última versión y levantar
docker compose -f docker-compose.aws.yml pull
docker compose -f docker-compose.aws.yml up -d
```

### ⏪ Rollback

Si una versión falla, volvés al SHA anterior:

```bash
IMAGE_TAG=sha-anterior docker compose -f docker-compose.aws.yml up -d
```

### 📦 Servicios

| Servicio | Imagen | Expuesto |
|---|---|---|
| `db` | `postgres:16-alpine` (Docker Hub) | Solo red interna |
| `api` | ECR `my-pets-api` | Solo red interna |
| `web` | ECR `my-pets-web` | Puerto `80` en el host |

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
| [`docker-compose.aws.yml`](docker-compose.aws.yml) | Deploy EC2 con ECR |
| [`scripts/build-and-push.sh`](scripts/build-and-push.sh) | Build + push local a ECR |
| [`.github/workflows/deploy-ecr.yml`](.github/workflows/deploy-ecr.yml) | CI/CD automático a ECR |

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
