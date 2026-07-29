# Working Rules

## Architecture
- Handwritten backend code lives under `internal/<domain>`; contracts under `proto/<domain>` (no `/v1` suffix — a deliberate v6 deviation from v5: this is a single-repo, single-deploy template with no external API consumers needing parallel versions, so path versioning is ceremony; see `docs/plan-v6.md` Issue 9).
- Register each domain once through `cmd/server/register_<domain>.go`; never add a second registry.
- SQL lives under `db/queries`; schema history lives under `db/migrations` (embedded, applied at server startup, up only, goose).
- Never hand-edit `internal/gen`, `web/src/lib/gen`, or `web/src/routeTree.gen.ts`; regenerate with `make gen` (routes regenerate on `vite dev`/`vite build` via the TanStack Router plugin).

## Security
- All Connect procedures are protected by default; public procedures require the explicit allowlist in `cmd/server/main.go`.
- The token only moves through `web/src/lib/auth.ts`.
- Never log OTPs, JWTs, database URLs, or secrets.
- Never create, drop, or reset PostgreSQL implicitly.
- Sensitive behavior (auth, authorization, money, deletion, destructive migration) requires owner approval before implementation.

## Docs
- Before developing UI, if `docs/route-catalog.html` or `docs/screenshots/` are missing from the working tree, run `make docs` first and use the output as the current-state reference — don't trust a stale mental model.
- Every new or changed route or visual state needs its fixture (`web/docs/fixtures/<route>.json`) updated in the same PR; `make docs` fails hard on a route without one.
- `docs/routes.json` is the one tracked, CI-gated artifact; `route-catalog.html` and `screenshots/` are gitignored and regenerable — never hand-author them.
- Each visually-distinct data regime gets its own named fixture state with realistic mock data (not a generic empty/loading/error enum).

## Decisions
- Read `docs/decisions/` before touching architecture.
- An `accepted` ADR is never edited — supersede it with a new one instead.
- Write a new ADR only for: swapping a major dependency, changing an already-in-use proto contract, changing the auth model, or changing the deployment shape.
- Do not write an ADR for features, styling, or anything revertable in one PR — that's just a PR.

## Testing
- No browser E2E; the `make docs` route-catalog fixtures are the integration surface.
- Go unit tests use a fake repo and an injected clock where time matters.
- While iterating, prefer targeted `go build`/`go test` (or `bun run typecheck`/`vitest`) on the touched packages; run the full `make check` once, as the final gate, before declaring work complete.

## Tooling
- Routine Bun/dependency bumps don't need an ADR — just a PR that passes `make check`.
