# ADR-001: React + TanStack Router (SPA) as the frontend stack

Status: accepted

Context: v5 used Svelte 5/SvelteKit; rebuilding the template from scratch is the point to reconsider the frontend, since swapping it later is expensive. The Go backend (Connect RPC, single binary with an embedded SPA) is unchanged.

Decision: Use React 19 + TanStack Router in SPA mode (file-based routes, no SSR) + TanStack Query via `@connectrpc/connect-query`, styled with Tailwind v4 + shadcn/ui.

Consequences:
- Generated Connect Query hooks replace v5's hand-written `client.ts` wrappers; caching/invalidation live in TanStack Query.
- shadcn/ui components live in `web/src/lib/components/ui/`, excluded from strict lint like other generated code.

Rejected alternatives:
- Continue with Svelte 5/SvelteKit: smaller/faster, but React has the deeper ecosystem and AI-agent familiarity for a fresh template.
- TanStack Start (SSR): adds a server-rendering runtime the single-binary Go backend doesn't need.
