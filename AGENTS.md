# Working Rules

## Execution
- Think Before Coding: inspect relevant code and one analog.
- Simplicity First: choose the smallest complete solution.
- Surgical Changes: touch only required files.
- Goal-Driven Execution: implement, verify, and finish.

## Project
- Backend: `internal`; contracts: `proto`; SQL: `db/queries`; migrations: `db/migrations`.
- Use `user` for backend patterns and `profile` for form patterns.
- Work in order: SQL, proto, backend, registration, route, fixture, generation, tests.
- Register domains once in `cmd/server/register_<domain>.go` and `main.go`.
- Never edit generated files. Batch SQL and proto before `make gen`.

## Safety
- Connect requires auth unless allowlisted in `cmd/server/main.go`.
- Keep tokens in `web/src/lib/auth.ts`. Never log secrets.
- Never create, drop, or reset PostgreSQL.
- Ask before auth, money, deletion, or destructive migration changes.

## Verify
- Add realistic fixtures for changed routes and states.
- Do not add feature ADRs or browser E2E tests.
- Use fake Go repositories. Inject clocks when needed.
- Run targeted tests while editing.
- For UI changes, run `make docs` once.
- Run `make check` once, last.
