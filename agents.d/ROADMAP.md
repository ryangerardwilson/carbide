# Roadmap

This roadmap is for agents and maintainers. Keep the public README focused on
human evaluation and first use.

## Current Baseline

- Source install from `cli/install.sh`.
- Compiled Go CLI with `new`, `init`, `project migrate`, `run dev`, `status`,
  `stop dev`, `follow logs`, `logs`, `doctor`, `doctor env`,
  `doctor runtime`, `doctor framework`, `deploy preview`, and guarded
  `deploy apply`.
- Generated `web`, `api`, and `db` services with Docker Compose watch.
- Bun/React/Tailwind browser app with register, login, logout, dashboard, and
  left-sidebar app shell.
- Go API backed by Postgres users and sessions.
- Environment/secrets contract in `carbide.toml`.
- Generated `AGENTS.md` and `agents.d` operating notes.
- Queryable structured dev logs in `.carbide/log/dev.jsonl`.
- Fast project doctor and Docker-backed runtime doctor.
- CI coverage for shell syntax, Go CLI tests, repo contract, scaffold checks,
  and generated Docker smoke flow.

## Phase 1: HTTP Core

- Harden routing for common HTTP methods beyond the starter routes.
- Add request parsing helpers for headers, query params, path params, and
  forms.
- Add response helpers for text, JSON, redirects, files, and errors.
- Add middleware chaining with predictable ownership rules.
- Add structured error pages for development and safe production errors.

## Phase 2: Application Kernel

- Harden the generated `web/`, `api`, and `db` directory contract.
- Expand configuration loading from environment and checked-in defaults.
- Harden environment contract validation and protected framework-owned keys.
- Add service registration without hidden reflection.
- Add logging with request IDs.
- Add graceful shutdown and worker lifecycle hooks.

## Phase 3: Frontend And Assets

- Keep the Bun/React web container as the public local-development entrypoint.
- Proxy `/api` and `/health` to the Go API service to preserve same-origin
  cookies.
- Keep the generated React starter useful without turning Carbide into a
  frontend package ecosystem.
- Make Tailwind the mandatory generated styling path.
- Serve the React shell with content-hashed JS and CSS assets by default.

## Phase 4: Database Layer

- Keep Postgres as the required database.
- Harden connection pooling.
- Add migrations with up/down support.
- Add a query builder with parameter binding by default.
- Add schema inspection helpers for Postgres-specific capabilities.

## Phase 5: Web App Essentials

- Harden register, login, logout, and dashboard.
- Add signed cookies and encrypted session storage.
- Add CSRF protection.
- Add validation primitives.
- Replace the starter password hash with a production-grade password hashing
  contract.
- Add file upload handling with size and type controls.

## Phase 6: Background Work

- Add Postgres-backed queues.
- Add scheduled jobs.
- Add mail driver contracts.
- Add cache contracts.
- Add retries, dead-letter behavior, and job inspection commands.

## Phase 7: Developer Experience

- Harden project scaffolding.
- Add migration generation.
- Add infrastructure generation, validation, and diff commands.
- Add test helpers for HTTP requests and database state.
- Add containerized watch/rebuild workflow.
- Add debug tooling for request lifecycle and connection leaks.

## Phase 8: Production Contract

- Define the official production image.
- Define the first production infrastructure-as-code target after local Compose
  is stable.
- Add health checks and readiness checks.
- Extend structured dev logs into the production container contract.
- Harden single-VM deploy apply beyond the docs app target.
- Implement clustered apply for previewed multi-VM environment targets.
- Add backup, restore, and migration rollback guidance.

## Phase 9: Ecosystem

- Stabilize extension points.
- Add first-party packages only where the core framework has repeated evidence.
- Document compatibility rules.
- Publish upgrade guides between framework versions.

