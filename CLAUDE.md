# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Conventions

- **Git commit messages must be in English** — always write commit messages in English, regardless of the language used in comments or documentation.

## Commands

All commands run from the **project root** via the root `Makefile`.

### Lifecycle

```bash
make setup    # First run: start Docker stack + apply schema + load seed data
make up       # Start stack — user data is preserved between restarts
make down     # Stop containers — data persists (volume is kept)
make clean    # Stop + delete volumes — all data is lost
```

### Rebuild (after code changes)

**IMPORTANT:** `make up` uses Docker cache and may NOT pick up code changes. After modifying source files, always use `make rebuild` to force a clean build:

```bash
make rebuild            # Rebuild frontend + API without cache, restart containers
make rebuild-frontend   # Rebuild only frontend (after changing frontend/src/)
make rebuild-api        # Rebuild only API (after changing backend/)
```

**When to use what:**
- `make up` — starting stopped containers, no code changes
- `make rebuild` — after ANY code change (frontend or backend)
- `make rebuild-frontend` — after frontend-only changes
- `make rebuild-api` — after backend-only changes

### Database

```bash
make migrate  # Apply schema only (001_create_tables.sql) — safe, idempotent, never touches data
make seed     # DESTRUCTIVE: truncate all tables and reload test data (002_seed_data.sql)
```

### Backend tests (no database required — uses sqlmock)

```bash
# From project root
make test

# From backend/ directory
go test ./... -v -count=1

# Single package
go test ./internal/repository/... -v
go test ./internal/handlers/... -v

# Single test function
go test ./internal/repository/... -run TestCompanyRepository_GetAll -v

# With coverage report
make test-coverage          # generates backend/coverage.html
make test-docker            # run tests inside an isolated Docker container
```

### Monitoring

```bash
make ps       # Container status
make logs     # All container logs (follow)
make logs-api # API logs only
```

### Local development (optional — without Docker for frontend/backend)

```bash
make dev-local      # starts MySQL in Docker, prints instructions
make migrate        # apply schema after MySQL is ready
make dev-api        # run Go API locally on :8081
make dev-frontend   # build + serve frontend locally on :3000 (Caddy)
```

## Architecture

### Request flow

```
Browser → Caddy (:3000) → React SPA
Browser → Go API (:8081) → Handler → DTO ↔ Model → Repository → MySQL (:3307 on host)
```

Data conversion pipeline:

```
Incoming:  JSON request → dto.JobRequest   → (ToModel())            → models.Job → Repository → MySQL
Outgoing:  MySQL        → Repository       → models.Job             → (JobResponseFromModel()) → dto.JobResponse → JSON response
```

The frontend API base URL is hardcoded in `frontend/src/services/api.ts`:
```typescript
const API_BASE_URL = 'http://localhost:8081/api/v1';
```

### Backend (`backend/`)

Layered Go architecture: **Handler → DTO → Repository → MySQL**.

- `cmd/api/main.go` — entry point: wires repositories, handlers, router, starts server
- `internal/database/db.go` — singleton `*sql.DB`, configured from env vars (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`)
- `internal/models/` — internal domain structs; nullable DB columns use `sql.NullString`, `sql.NullFloat64`, etc. Models with nullable fields have `json:"-"` tags — they are **never serialized directly**; all JSON conversion goes through DTO
- `internal/dto/` — Data Transfer Objects for API layer. Separate Request/Response structs per entity with mapper functions. Handles conversion between `sql.Null*` types (model) and pointer types `*string`, `*float64` (JSON). This layer is the **only place** where JSON shape is defined for entities with nullable fields
- `internal/repository/` — one struct per entity wrapping `*sql.DB` with raw SQL; `interfaces.go` defines `JobRepositoryInterface`, `CompanyRepositoryInterface`, `JobSkillRepositoryInterface` used by handlers
- `internal/handlers/` — one struct per entity; accepts repository interfaces (not concrete types), enabling mock-based testing. Handlers decode incoming JSON into `dto.*Request`, call `ToModel()` to get a domain model, pass it to the repository, then convert the result back via `dto.*ResponseFromModel()`

#### DTO conventions

Each entity with nullable fields gets a file in `internal/dto/` containing:

| Type | Purpose | Example |
|---|---|---|
| `*Request` | Incoming JSON → model | `dto.JobRequest` — decoded from request body, `ToModel()` returns `models.Job` |
| `*Response` | Model → outgoing JSON | `dto.JobResponse` — all nullable fields are `*string` / `*float64` (serialize to value or `null`) |
| `*ResponseFromModel()` | Single model mapper | `dto.JobResponseFromModel(j models.Job) JobResponse` |
| `*ResponseList()` | Slice mapper | `dto.JobResponseList(jobs []models.Job) []JobResponse` |

Entities without nullable fields (`Company`, `Skill`, `JobSkill`, stats models) still use `json` tags directly on the model struct — no DTO needed until the API shape diverges from the DB shape.

**When to add a new DTO:** if the model uses `sql.Null*` types, or if the API response shape should differ from the internal model (e.g., embedding related data, hiding fields).

**CORS:** `corsMiddleware` wraps the entire `http.Server.Handler` — NOT registered via `r.Use()`. This is intentional: gorilla/mux returns 405 before middleware fires for unmatched methods, so wrapping the handler is the only way to handle OPTIONS preflight correctly.

**StatsHandler** uses `*repository.StatsRepository` directly (no interface) because stats queries are read-only and not tested with mocks.

### Testing strategy

Repository tests use `github.com/DATA-DOG/go-sqlmock v1.5.2` to mock `*sql.DB`. Handler tests use `net/http/httptest` with in-package mock structs (`mockJobRepo`, `mockCompanyRepo`) that implement the repository interfaces. No real database is needed to run tests.

Handler tests send `dto.*Request` structs as request bodies and decode responses into `dto.*Response` structs — never into `models.*` directly.

`AddRow` in go-sqlmock v1 requires `[]driver.Value`, not `[]interface{}` — numeric values must be cast to `int64`.

### Frontend (`frontend/`)

React 18 SPA bundled with esbuild (no Webpack/Vite). State is managed with MobX in `src/stores/RootStore.ts`. Charts use Recharts. Routing uses React Router v6.

**Build pipeline:** esbuild compiles `src/` → `dist/bundle.js`. Caddy serves `dist/` as a SPA with `try_files {path} /index.html`. The Docker image runs a multi-stage build (node → caddy). Code changes require `make rebuild-frontend` to take effect.

### Docker

- `docker-compose.full.yml` (root) — full stack: MySQL + API + Frontend on a shared bridge network `job_stats_network`
- `backend/docker-compose.yml` — backend-only (MySQL + API)
- `backend/docker-compose.test.yml` — isolated test runner; mounts `backend/coverage/` for the coverage report
- MySQL data is stored in the named volume `mysql_data`; `make down` preserves it, `make clean` removes it

**Docker cache caveat:** `make up` runs `docker-compose up -d --build`, but Docker may cache layers if source files haven't changed in a way that invalidates the COPY layer. Use `make rebuild` (which uses `--no-cache`) when `make up` doesn't pick up changes.

### Migrations

- `backend/migrations/001_create_tables.sql` — schema with `CREATE TABLE IF NOT EXISTS`; safe to run repeatedly
- `backend/migrations/002_seed_data.sql` — starts with `TRUNCATE` (via `SET FOREIGN_KEY_CHECKS=0`); always destructive
- `migrate.sh` runs only schema; `seed.sh` runs only seed data
- Both scripts detect the MySQL container by name (`job_stats_mysql`) using `docker ps`, not `docker-compose ps`, so they work regardless of which compose file started the container

## Future: DDD migration path

The current DTO architecture is designed as a stepping stone toward Domain-Driven Design. Planned evolution:

```
Current:   Handler → DTO ↔ Model (sql.Null*) → Repository → MySQL
Future:    Handler → DTO ↔ Domain Entity (pure Go) → Repository (sql.Null*) → MySQL
```

When migrating to DDD:
1. Extract pure domain entities (no `sql.Null*`, no `json` tags) into `internal/domain/`
2. Move `sql.Null*` handling into repository-layer mappers
3. DTO layer stays unchanged — it already works with clean Go types (`*string`, `*float64`)
4. Add domain services for business logic that currently lives in handlers
