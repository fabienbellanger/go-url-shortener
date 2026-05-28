# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository layout

This is a two-part monorepo for a URL shortener:

- `server/` — Go backend (Fiber + GORM + MySQL), exposing a JSON API under `/api/v1` and serving short-link redirects at `/:id`. Entry point: `server/cmd/main.go` → `cli.Execute()`.
- `client/` — Quasar/Vue 3 + TypeScript admin SPA that consumes the API.

The server is invoked via Cobra subcommands (no plain "start the binary"): `run`, `register`, `logs`. Configuration is loaded from `server/.env` via Viper (`viper.SetConfigFile(".env")` in `server/cli/root.go`); this file is required and not in the repo, so a missing `.env` is the typical first-run failure.

## Server (Go) — commands

All from `server/`:

```bash
make serve          # go run cmd/main.go run
make serve-race     # with -race
make watch          # air live-reload
make build          # produces ./go-url-shortener binary
make test           # go test -cover ./...
make test-verbose   # adds -v
make bench          # go test -benchmem -bench=. ./...
make logs           # tails server logs via the CLI's `logs --server`
make update         # go get -u ./... && go mod tidy
```

Run a single test: `cd server && go test ./utils -run TestName -v`.

Server CLI (after `make build`, or via `go run cmd/main.go <cmd>`):
- `go-url-shortener run` — start the HTTP server
- `go-url-shortener register -l <last> -f <first> -e <email> -p <pass>` — create an admin user (password ≥ 8 chars). This is the only way to seed the first user.
- `go-url-shortener logs --server` / `--database` — colorized log readers for the zap/GORM logs

Coverage HTML: `make view-cover-count` (count mode) or `make view-cover-atomic`.

Docker stack (server + MySQL 5.7): `cd server && docker-compose up` exposes the API on `localhost:9900`.

## Client (Quasar) — commands

From `client/`:

```bash
npm install
npm run dev         # quasar dev (Vite, hot reload)
npm run build       # quasar build (production SPA in dist/)
npm run lint        # eslint .js,.ts,.vue
npm run format      # prettier write
```

There is no test suite on the client (`npm test` is a no-op).

Quasar config: copy `client/quasar.config.mjs.dist` to `client/quasar.config.mjs` and adjust the API base URL before `npm run dev` (the commit `9ee9bb7` was specifically about fixing this).

## Architecture notes

### Server request flow

`server/server.go::Run` builds a Fiber app, then wires routes through `server/routes.go`:

1. **Public web routes** (`registerPublicWebRoutes`): the short-link redirector `GET /:id` (handled by `handlers.RedirectURL`), the Basic-Auth-protected `/doc/api-v1`, and the `/assets` static filesystem.
2. **Public API routes** (`registerPublicAPIRoutes`, prefix `/api/v1`): `POST /login`, `POST /forgotten-password/:email`, `PATCH /update-password/:token`.
3. **JWT gate** (`initJWT`) — everything registered after this requires a Bearer token.
4. **Protected API routes** (`registerProtectedAPIRoutes`): user CRUD under `/api/v1/users`, link CRUD + CSV import/export under `/api/v1/links`.

Order matters: the catch-all 404 handler is mounted after JWT, so adding new public routes must happen before the `initJWT(app)` call in `Run`.

### Layered packages

Standard repository → handler split:

- `models/` — GORM model structs (`User`, `Link`, `PasswordResets`). Listed in `db/migration.go::modelsList` for `AutoMigrate`; new models must be added there.
- `repositories/` — DB access functions taking a `*db.DB` and returning models/errors. No Fiber types here.
- `handlers/` — Fiber HTTP handlers that validate input (via `utils.ValidateStruct`, go-playground/validator), call repositories, and return JSON / `utils.HTTPError`.
- `db/` — connection setup (`New`), GORM logger config (output can be file or stdout, level driven by `APP_ENV` + `GORM_LOG_LEVEL`), and the migration list.
- `utils/` — short-code generation (`url_shortener.go`, base58-go), validator wrapper, HTTP error shape.
- `cli/` — Cobra commands. `initConfigLoggerDatabase(initLogger, initDatabase bool)` is the shared bootstrap; pass `false` for either when a subcommand doesn't need it (e.g. `register` skips the logger).

### Auto-migrations

When `DB_USE_AUTOMIGRATIONS=true`, `cli/server.go::startServer` calls `db.MakeMigrations()` before `server.Run`. This runs GORM `AutoMigrate` on every model in `db/migration.go::modelsList` — schema changes go through model struct edits, not hand-written SQL migrations.

### Config keys (Viper / `.env`)

Read directly via `viper.GetString/Bool/Int/Duration` throughout the code (no central config struct). Notable groups:
- `APP_NAME`, `APP_ENV`, `APP_ADDR`, `APP_PORT`
- `SERVER_PREFORK`, `SERVER_BASICAUTH_USERNAME`, `SERVER_BASICAUTH_PASSWORD`
- `DB_DRIVER`, `DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`, `DB_DATABASE`, `DB_CHARSET`, `DB_COLLATION`, `DB_LOCATION`, `DB_MAX_IDLE_CONNS`, `DB_MAX_OPEN_CONNS`, `DB_CONN_MAX_LIFETIME` (hours), `DB_USE_AUTOMIGRATIONS`
- `GORM_LOG_LEVEL`, `GORM_LOG_OUTPUT`, `GORM_LOG_FILE_PATH`

JWT signing key, expiration, and other middleware settings are similarly pulled from Viper inside `server.go` / `routes.go`.

### Templates & assets

Fiber view engine points at `server/templates/*.gohtml` (used by the API doc page). Static files live in `server/assets/` and are served under `/assets` with a 1-hour cache.

## Routes reference

`server/ROUTES.md` documents every endpoint with example payloads; `server/ROUTES.http` has ready-to-fire requests for VS Code REST Client / JetBrains HTTP client.

## Benchmarking

`server/drill.yml` is a [Drill](https://github.com/fcsonline/drill) scenario: `drill --benchmark drill.yml --stats --quiet` from `server/`.
