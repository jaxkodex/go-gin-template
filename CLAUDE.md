# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build (requires generated API file — see OpenAPI section below)
make build

# Run
go run main.go

# Test
go test -v ./...

# Single test
go test -v ./src/... -run TestName

# Add/update dependencies
go get <package>
go mod tidy

# Database migrations (requires goose + a running Postgres)
make migrate-up
make migrate-down
make migrate-status
make migrate-create name=<migration_name>

# Generate OpenAPI server code from a remote spec
make generate-api url=https://example.com/api.json
```

## Architecture

`main.go` is the entry point. It loads config, connects to Postgres, initializes auth middleware, and wires routes before starting the Gin HTTP server.

### Package layout

| Package | Role |
|---------|------|
| `src/config` | Typed config struct populated from env vars via `config.Load()`. Single source of truth for all env-var reads (the database package pre-dates this and reads env vars directly). |
| `src/database` | Global `pgxpool.Pool` (`database.DB`) initialized by `database.Connect()`. |
| `src/middleware` | Auth orchestrator (`auth.go`) that tries Firebase then API key. Each verifier (`firebase_auth.go`, `api_key.go`) is independent and returns a result without aborting — the orchestrator decides. |
| `src/routes` | One file per route group. `RegisterRoutes` receives the auth middleware and creates a protected route group under `/`. Public routes (e.g. `/health`) are registered directly on the engine, outside the group. |
| `src/api` | Hand-written `server.go` (`Server` struct implementing the generated interface) + gitignored `api.gen.go` produced by `oapi-codegen`. |

### Authentication flow

`middleware.NewAuth(cfg)` initializes whichever verifiers are enabled and returns an `*AuthMiddleware`. `auth.Authenticate()` is a Gin handler factory:
- Both disabled → no-op pass-through (dev mode)
- Either enabled → tries Firebase Bearer token first, API key second; first success wins; all fail → 401 `{"error":"unauthorized","code":"AUTH_REQUIRED"}`

Downstream handlers retrieve identity via `middleware.GetAuthUserID`, `GetAuthMethod`, `GetAuthClaims`.

### OpenAPI code generation

`src/api/api.gen.go` is **gitignored** — never commit it. It is generated on demand:

```bash
make generate-api url=<remote-spec-url>
```

`src/api/server.go` contains the hand-written `Server` struct that implements the generated `ServerInterface`.

## Environment variables

Copy `.env.sample` to `.env`. Key variables:

| Variable | Purpose |
|----------|---------|
| `SERVER_PORT` | HTTP listen port (default `8000`) |
| `DB_*` | Postgres connection (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_MAX_CONNS`, `DB_MIN_CONNS`) |
| `AUTH_FIREBASE_ENABLED` | Enable Firebase auth verifier |
| `FIREBASE_CREDENTIALS_FILE` | Path to service account JSON |
| `FIREBASE_CREDENTIALS_JSON` | Inline service account JSON (alternative) |
| `AUTH_API_KEY_ENABLED` | Enable API key verifier |
| `AUTH_API_KEYS` | Comma-separated list of valid API keys |

## Git conventions

- **Never commit generated files.** `src/api/api.gen.go` is gitignored; regenerate it locally with `make generate-api`.
