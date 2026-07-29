# Template v6

An AI-first full-stack template with a **Go 1.26, Connect RPC, and PostgreSQL**
backend; a **React 19, TanStack Router (SPA), and shadcn/ui** frontend; and a
single production binary. Working rules for agents and contributors live in
[`AGENTS.md`](AGENTS.md) (symlinked as `CLAUDE.md`); architectural decisions
that are expensive to reverse are recorded in [`docs/decisions/`](docs/decisions/).

## Setup

Requirements: Go (version pinned in [`go.mod`](go.mod)), Bun (version pinned
as `packageManager` in [`web/package.json`](web/package.json)), and an
external PostgreSQL.

```bash
cp .env.example .env     # set DATABASE_URL and a >=32-byte JWT_SECRET
make bootstrap           # installs all dependencies (the only target that does)
make run                 # build the frontend, then run the single server process
```

Schema migrations are embedded in the binary and applied automatically at
server startup (up only); no target here creates, drops, or resets the
database itself.

## Make targets

| Target | Purpose |
|---|---|
| `make bootstrap` | Install/download all dependencies (Go modules, Bun packages, Playwright's Chromium). The only target that installs anything. |
| `make gen` | Regenerate code from proto (buf) and SQL (sqlc). |
| `make check` | Done-signal: codegen, buf lint, gofmt, go vet, go test, web lint, typecheck, vitest, web build, go build. Run this once as the final gate before calling work done. |
| `make run` | Build the frontend once, then run the single Go server process. |
| `make build` | Produce `bin/server` with the SPA embedded. |
| `make docs` | Regenerate `docs/routes.json` (tracked), `docs/route-catalog.html`, and `docs/screenshots/` (both gitignored, regenerable) — no backend/DB needed. |
| `make help` | List all commands. |

## Security notes

- The seeded `admin@localhost` account accepts the static OTP `123456`
  (exact-match only). Remove or protect it before any untrusted deployment.
- Other accounts use a 5-minute TOTP derived from `JWT_SECRET`; no email
  provider is wired up unless `MAIL_URL` is set.
- An OTP can be replayed within its own 5-minute step (documented limitation).
- The bearer token is stored in `localStorage`; an XSS in this origin could
  read it.
