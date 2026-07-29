# ADR-003: routes.json is the only tracked docs artifact

Status: accepted

Context: `make docs` produces `docs/routes.json`, `docs/route-catalog.html`, and `docs/screenshots/`, all fully regenerable from source. Committing all three risks a stale copy being trusted as current.

Decision: Track only `docs/routes.json` (small, diffable, CI-gated); gitignore `route-catalog.html` and `screenshots/`.

Consequences:
- CI fails when a route/fixture changes without a matching `routes.json` update, without storing binary PNGs in git.
- Reviewers must run `make docs` locally to see the HTML catalog or screenshots — no stale checked-in copy to mistakenly trust.
- `AGENTS.md` codifies running `make docs` before UI work if these files are missing, and again after UI changes.

Rejected alternatives:
- Commit all three artifacts: guarantees stale-diff bugs and bloats the repo with binary screenshot diffs.
- Track nothing, not even routes.json: loses the CI diff-gate for a route/state shipped without an updated fixture.
