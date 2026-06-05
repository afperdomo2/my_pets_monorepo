# 🧪 Testing — Backend API

> Guía completa de tipos de prueba, cobertura, ejecución y convenciones para el backend Go.

---

## 📦 Stack tecnológico

| Componente | Propósito | Enlaces |
|---|---|---|
| [`testing`](https://pkg.go.dev/testing) | Runner estándar de Go | — |
| [`testify/mock`](https://github.com/stretchr/testify) | Creación de mocks para repositorios | `docs/api/conventions.md` |
| [`testify/require`](https://github.com/stretchr/testify) | Aserciones con fail inmediato | — |
| [`httptest`](https://pkg.go.dev/net/http/httptest) | `ResponseRecorder` + `gin.CreateTestContext` | — |
| [`testcontainers-go`](https://golang.testcontainers.org/) | Postgres 16 en contenedor Docker | `internal/database/testutil/testdb.go` |

---

## 🎯 Tipos de prueba

### 🔹 1. Handler tests — unitarios con mocks

Prueban la **capa HTTP** sin base de datos. Los repositorios se reemplazan con mocks.

| Característica | Detalle |
|---|---|
| ⚡ Velocidad | ~0.01s por test |
| 📁 Ubicación | `internal/domain/<name>/handler_test.go` |
| 🧩 Dependencias | `testify/mock`, `httptest`, `gin.CreateTestContext` |
| 🐳 ¿Requiere Docker? | ❌ No |
| 🔄 Cacheable | ✅ Sí |

**Escenarios típicos:** `200 success`, `4xx` (bad request, forbidden, not found, conflict), `5xx` (error de repositorio), validación de payloads, permisos (system vs non-system), paginación por defecto, límite de mascotas alcanzado

---

### 🔹 2. Pure unit tests

Lógica de negocio pura, sin mocks ni base de datos.

| Característica | Detalle |
|---|---|
| ⚡ Velocidad | ~0.001s por test |
| 📁 Ubicación | `internal/domain/<name>/xxx_test.go` |
| 🧩 Dependencias | `testify/require` |
| 🐳 ¿Requiere Docker? | ❌ No |
| 🔄 Cacheable | ✅ Sí |

**Escenarios típicos:** cálculo de etapas de vida (`[pet/lifestage_test.go](/apps/api/internal/domain/pet/lifestage_test.go)`), validaciones, helpers de fechas

**Patrón:** table-driven con `t.Run()` por subtest

---

### 🔹 3. Repository tests — integración con Postgres real

Prueban los **queries reales** contra PostgreSQL usando `testcontainers-go`.

| Característica | Detalle |
|---|---|
| ⚡ Velocidad | ~3s por test (incluye levantar container) |
| 📁 Ubicación | `internal/domain/<name>/gorm_repo_test.go` |
| 🧩 Dependencias | `testcontainers-go`, `database/testutil` |
| 🐳 ¿Requiere Docker? | ✅ Sí — indispensable |
| 🔄 Cacheable | ⚠️ Sí, pero se puede forzar con `-count=1` |

**¿Cómo funciona?**

```
  1. testcontainers llama a la Docker API
  2. Descarga postgres:16-alpine (si no está en caché local)
  3. Crea el contenedor con credenciales test
  4. Ejecuta AutoMigrate (mismas migraciones que la app real)
  5. Se insertan datos de prueba (seed)
  6. Corre el método del repositorio
  7. Se verifican los resultados
  8. Se destruye el contenedor automáticamente
```

**Escenarios típicos:** queries del dashboard (pending_tasks, overdue_tasks), aislamiento entre usuarios por `ownerID`, paginación real, constraints de BD, foreign keys y cascade, owner-scoped CRUD

> ⚠️ **Caché de Go:** si el código no cambió, Go reusa el resultado anterior. Los repo tests se marcan como `(cached)` y no ejecutan containers.
>
> Para forzar ejecución:
> ```bash
> go test -count=1 ./internal/domain/dashboard/
> go clean -testcache
> ```

---

## 📊 Cobertura actual

| Dominio | 🧪 Handler | 🧬 Pure unit | 🗄️ Repo | 📊 Total |
|---|---|---|---|---|
| [`auth`](/apps/api/internal/domain/auth/) | 11 | 7 | — | **18** |
| [`dashboard`](/apps/api/internal/domain/dashboard/) | 2 | — | 4 | **6** |
| [`exam`](/apps/api/internal/domain/exam/) | 17 | — | 8 | **25** |
| [`health`](/apps/api/internal/domain/health/) | 1 | — | — | **1** |
| [`health_catalog`](/apps/api/internal/domain/health_catalog/) | 22 | — | 7 | **29** |
| [`health_record`](/apps/api/internal/domain/health_record/) | 19 | — | 8 | **27** |
| [`pet`](/apps/api/internal/domain/pet/) | 11 | 27 | 8 | **46** |
| [`setup`](/apps/api/internal/domain/setup/) | 8 | — | — | **8** |
| [`user`](/apps/api/internal/domain/user/) | 20 | — | 11 | **31** |
| [`vaccine_application`](/apps/api/internal/domain/vaccine_application/) | 19 | — | 8 | **27** |
| [`validation`](/apps/api/internal/validation/) | — | 9 | — | **9** |
| **Total** | **130** | **43** | **54** | **227** |

```text
📈 Cobertura por tipo:
  ████████████████████████████████████████████████████████████████████░░░░  Handler  (130 tests)
  █████████████████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  Pure unit (43 tests)
  ██████████████████████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  Repo     (54 tests)
```

### 🗺️ Dominios con cobertura completa

| Dominio | Handler | Repo |
|---|---|---|
| `auth` | ✅ | ❌ |
| `dashboard` | ✅ | ✅ |
| `exam` | ✅ | ✅ |
| `health` | ✅ | ❌ |
| `health_catalog` | ✅ | ✅ |
| `health_record` | ✅ | ✅ |
| `pet` | ✅ | ✅ |
| `setup` | ✅ | ❌ |
| `user` | ✅ | ✅ |
| `vaccine_application` | ✅ | ✅ |

---

## 🚀 Cómo ejecutar

### Todos los tests

```bash
make test-api
```

### Por dominio

```bash
go test ./internal/domain/auth/              -v
go test ./internal/domain/dashboard/         -v
go test ./internal/domain/exam/              -v
go test ./internal/domain/health/            -v
go test ./internal/domain/health_catalog/    -v
go test ./internal/domain/health_record/     -v
go test ./internal/domain/pet/               -v
go test ./internal/domain/setup/             -v
go test ./internal/domain/user/              -v
go test ./internal/domain/vaccine_application/ -v
```

### Por tipo de test

```bash
# 🧪 Solo handler tests (rápidos, sin DB)
go test ./internal/domain/pet/ -run 'Test(PetHandler|GetSummary)' -v

# 🗄️ Solo repository tests (con Postgres, requieren Docker)
go test ./internal/domain/pet/ -run 'TestGormRepo' -v

# 🧬 Solo pure unit tests
go test ./internal/domain/pet/ -run 'TestCalculate' -v
```

### Sin caché (forzar ejecución real)

```bash
# Forzar ejecución (útil para repo tests)
go test ./internal/domain/dashboard/ -count=1 -v

# Borrar todo el caché de tests
go clean -testcache
```

---

## 📐 Convenciones

### 📁 Nomenclatura de archivos

| Sufijo | Contenido | Ejemplo |
|---|---|---|
| `handler_test.go` | Handler tests con mocks | `dashboard/handler_test.go` |
| `gorm_repo_test.go` | Repository tests con testcontainers | `dashboard/gorm_repo_test.go` |
| `lifestage_test.go` | Pure unit: cálculos de etapa de vida | `pet/lifestage_test.go` |
| `xxx_test.go` | Otros pure unit tests | — |

### 🧪 Patrón de handler tests

```go
type mockRepo struct{ mock.Mock }

func (m *mockRepo) Metodo(ctx context.Context, args...) (Resultado, error) {
    a := m.Called(ctx, args...)
    return a.Get(0).(Resultado), a.Error(1)
}

func TestHandler_Escenario(t *testing.T) {
    m := new(mockRepo)
    h := NewHandler(m)
    m.On("Metodo", mock.Anything, "arg").Return(resultado, nil)

    c, w := setupGin()
    c.Set("userID", "user-1")
    handler.GetX(c)

    require.Equal(t, 200, w.Code)
    m.AssertExpectations(t)
}
```

```go
// Helper común para crear contexto Gin de prueba
func setupGin() (*gin.Context, *httptest.ResponseRecorder) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("GET", "/", nil)
    return c, w
}
```

### 🗄️ Patrón de repository tests

```go
func TestGormRepo_Metodo(t *testing.T) {
    db, cleanup := testutil.NewPostgres(t) // ← levanta Postgres + migraciones
    defer cleanup()

    seedData(t, db)                        // ← inserta datos de prueba

    repo := NewGormRepo(db)
    resultado, err := repo.Metodo(ctx, args...)

    require.NoError(t, err)
    require.Equal(t, valorEsperado, resultado)
}
```

### 🌱 Seed de datos — buenas prácticas

- **IDs**: usar UUIDs válidos — formato `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`
  - Helper recomendado: `uid(n)` → `00000000-0000-0000-0000-0000000000XX`
- **Fechas**: usar `time.Now().UTC()` con `AddDate` — así los tests **nunca expiran**
- **Aislamiento**: cada test crea sus propios datos (no comparten estado entre tests)
- **Limpieza**: el cleanup del container se encarga de todo — no hay state residual

---

## 📚 Referencias

| Recurso | Descripción |
|---|---|
| [`docs/api/conventions.md`](conventions.md) | Convenciones generales del backend Go |
| [`docs/api/commands.md`](commands.md) | Comandos de desarrollo y testing |
| [`internal/database/testutil/testdb.go`](/apps/api/internal/database/testutil/testdb.go) | Helper para levantar Postgres via testcontainers |
| [`Makefile`](/Makefile) | Targets `test-api`, `lint-api`, `tidy` |

## 🔍 Notas técnicas

### ¿Por qué testcontainers y no SQLite?

Los queries del proyecto usan funciones específicas de PostgreSQL:

| Característica | PostgreSQL | SQLite |
|---|---|---|
| `CURRENT_DATE` | ✅ | ❌ |
| `gen_random_uuid()` | ✅ | ❌ |
| Tipo `uuid` | ✅ | ❌ |
| `type:date` | ✅ | ❌ |

SQLite no sería compatible, por lo que testcontainers con Postgres real es la opción correcta.

### ¿Por qué no una DB compartida vía Docker Compose?

- Testcontainers **aisla cada test**: no hay estado residual entre ejecuciones
- No requiere configuración previa ni scripts de setup
- Cada suite de tests arranca con datos limpios y predecibles
- Es el estándar en la comunidad Go para integration tests
