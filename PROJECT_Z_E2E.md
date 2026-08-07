# Project Z E2E — Repository Capability Summary

This document is a factual capability inventory derived from committed repository evidence. It is documentation only and does not alter application code, dependencies, workflows, configuration, or runtime behavior.

## Runtime & Framework

- **Backend runtime:** Go, version `1.26.5` (`go.mod`).
- **Backend RPC framework:** [Connect](https://connectrpc.com) (`connectrpc.com/connect v1.20.0`, `connectrpc.com/otelconnect`), generated from Protobuf contracts (`go.mod`, `proto/`).
- **Database driver/migrations:** `github.com/jackc/pgx/v5` (PostgreSQL driver) and `github.com/pressly/goose/v3` (migrations), embedded and applied by the server at startup (`go.mod`, `db/embed.go`, `Makefile`).
- **Codegen tools:** `buf` (proto) and `sqlc` (SQL), pinned as Go tool dependencies (`go.mod` `tool (...)` block).
- **Frontend runtime:** Bun `1.3.14`, pinned via `"packageManager": "bun@1.3.14"` in `web/package.json` and enforced by `Makefile`'s `_check-tools` target.
- **Frontend framework:** React `^19.2.8` with Vite `^8.1.5`, TanStack Router `^1.170.18`, TanStack Query `^5.101.4`, and Tailwind CSS `^4.3.3` (`web/package.json`).

## Package Manager(s)

- **Go modules** — `go.mod` / `go.sum` at the repository root.
- **Bun** — `web/package.json` (`packageManager: bun@1.3.14`) with committed lockfile `web/bun.lock`.

No other package manager manifests or lockfiles (npm, yarn, pnpm) are present under `web/`.

## Available Commands

Commands below are quoted exactly as defined in `Makefile` and `web/package.json`.

| Category | Command | Source |
|---|---|---|
| Install | `make bootstrap` (runs `go mod download`, `cd web && bun install --frozen-lockfile`, `bunx playwright install chromium`) | `Makefile` |
| Lint (Go) | `$(BUF) lint` (buf lint) and `gofmt -l $(GOFMT_PATHS)` | `Makefile` (`check` target) |
| Lint (web) | `cd web && bun run lint` → `prettier --check . && eslint .` | `Makefile`, `web/package.json` |
| Typecheck | `cd web && bun run typecheck` → `tsc --noEmit` | `Makefile`, `web/package.json` |
| Test (Go) | `go vet $(GO_PKGS)` and `go test $(GO_PKGS)` (`GO_PKGS := ./cmd/... ./internal/... ./db ./web`) | `Makefile` |
| Test (web) | `cd web && bun run test:unit -- --run --passWithNoTests` → `vitest` | `Makefile`, `web/package.json` |
| Build (web) | `cd web && bun run build` → `vite build` | `Makefile`, `web/package.json` |
| Build (Go) | `go build $(GO_PKGS)` (via `make check`) or `go build -o bin/server ./cmd/server` (via `make build`) | `Makefile` |
| Codegen | `make gen` → `$(BUF) generate` and `$(SQLC) generate` | `Makefile` |
| Combined verification | `make check` (runs codegen, buf lint, gofmt check, go vet, go test, web lint, web typecheck, web unit tests, web build, go build in sequence) | `Makefile` |
| Docs | `make docs` → `cd web && bun run docs` (`vite build && bun docs/generate.ts`), regenerates `docs/routes.json`, `docs/templates.json`, `docs/route-catalog.html`, `docs/screenshots/` | `Makefile`, `web/package.json` |
| Run | `make run` (build frontend once, then `go run ./cmd/server`) | `Makefile` |

No standalone `install`-only command exists for either ecosystem outside of `make bootstrap` (which installs both Go module downloads and Bun dependencies together, plus a Playwright browser). There is no repository-defined command category left undocumented above; all install/lint/typecheck/test/build categories have at least one defined command.

## CI

`.github/workflows/ci.yml` runs, on push to `main` and on pull requests: `make bootstrap`, `make check BUF=buf SQLC=sqlc`, `make docs`, and `git diff --exit-code docs/routes.json` (failing the build if generated route docs are out of date).

## Important Source & Test Directories

**Backend (Go):**
- `cmd/server` — server entrypoint and domain registration (`main.go`, `register_auth.go`, `register_frontend.go`, `register_user.go`).
- `internal/auth` — auth service, JWT, OTP, delivery, interceptor (with colocated `_test.go` files).
- `internal/user` — user service and repository (with colocated `_test.go` file).
- `internal/health` — health check handler (with colocated `_test.go` file).
- `internal/mail` — mail sending (with colocated `_test.go` file).
- `internal/template` — HTML email templates (with colocated `_test.go` file).
- `internal/platform` — `config`, `database`, `server` platform wiring.
- `db/queries` — SQL queries (`auth.sql`, `user.sql`); `db/migrations` — Goose SQL migrations (`00001_users.sql`, `00002_profile.sql`); `db/embed.go` — embeds migrations into the binary.
- `proto/auth`, `proto/user` — Protobuf API contracts.

**Frontend (web):**
- `web/src/routes` — TanStack Router route tree (`index.tsx`, `__root.tsx`, `_authenticated.tsx`, `_authenticated/profile.tsx`, `login/index.tsx`, `login/otp.tsx`).
- `web/src/lib` — shared client/auth logic, including colocated test files `auth.test.ts` and `login.test.ts`.
- `web/src/lib/components` — shared UI components.

Go tests found: `internal/auth/{delivery,interceptor,jwt,otp,service}_test.go`, `internal/health/handler_test.go`, `internal/mail/mail_test.go`, `internal/template/template_test.go`, `internal/user/service_test.go`.

TypeScript tests found: `web/src/lib/auth.test.ts`, `web/src/lib/login.test.ts`.

## Project Z Template Files

No Project Z template files (e.g., any file matching `project_z*`, `projectz*`, or a pre-existing `PROJECT_Z_E2E.md`) were found anywhere in the repository. A repository-wide filename search returned no matches. Project Z template files are **absent**.
