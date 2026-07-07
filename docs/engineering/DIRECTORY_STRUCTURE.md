# Directory Structure

This page describes the project layout users get after running
`carbide new "My Carbide App"` or `carbide init`. Carbide's own repository
layout is contributor material; the product docs should lead with the generated
app because that is the structure application teams live in.

```text
my-carbide-app/
|-- .env.example
|-- .gitignore
|-- carbide.toml
|-- docker-compose.yml
|-- api/
|   |-- Dockerfile
|   |-- auth.go
|   |-- go.mod
|   |-- go.sum
|   |-- main.go
|   `-- routes.go
|-- db/
|   |-- go.mod
|   |-- go.sum
|   |-- migration/
|   |   `-- 001_auth.sql
|   |-- session.go
|   `-- user.go
`-- web/
    |-- Dockerfile
    |-- bun.lock
    |-- index.html
    |-- package.json
    |-- tsconfig.json
    `-- src/
        |-- component/
        |   |-- l1/
        |   |   |-- Button.tsx
        |   |   |-- Field.tsx
        |   |   |-- Surface.tsx
        |   |   |-- Text.tsx
        |   |   |-- ThemeToggle.tsx
        |   |   |-- index.ts
        |   |   `-- tokens.ts
        |   |-- l2/
        |   |   |-- AuthForm.tsx
        |   |   |-- Layouts.tsx
        |   |   `-- index.ts
        |   `-- l3/
        |       |-- AuthView.tsx
        |       |-- DashboardView.tsx
        |       |-- LoadingView.tsx
        |       `-- index.ts
        |-- lib/
        |   |-- cx.ts
        |   `-- types.ts
        |-- main.tsx
        |-- server.ts
        |-- styles.d.ts
        |-- styles.css
        `-- write-index.ts
```

## Root Contract

The generated project root is intentionally small:

- `web/`, `api/`, and `db/` map to standalone Docker services.
- `docker-compose.yml` owns local runtime coordination across services.
- `carbide.toml` owns the app name, slug, default dev port, product contract
  current starter runtime defaults, environment contract, and deploy targets.
- `.env.example` documents local development variables without storing real
  secrets.
- Carbide does not scaffold `README.md` or `AGENTS.md`. If the app owner
  creates them later, they are local prose files rather than framework-owned
  contract files.

There is no root `src/`, `frontend/`, `backend/`, `model/`, `controller/`,
`view/`, `infra/`, or `doc/` directory in generated apps.

## Service Directories

- `web/`: Bun, React, and Tailwind. This is the public browser entrypoint.
  It serves the browser app, proxies `/api` and `/health` to the API service,
  and owns content-hashed browser assets.
- `api/`: Go HTTP/API service. It owns auth, sessions, validation, routing,
  JSON responses, and the API Dockerfile.
- `db/`: Postgres-facing data module. It owns data access helpers and checked-in
  migration state.

## Frontend Structure

The generated web app uses Tailwind and keeps component tiers explicit:

- `web/src/component/l1/`: primitives and Tailwind utility tokens, including
  the built-in light/dark scrollbar utility group.
- `web/src/component/l2/`: reusable composed patterns such as forms and
  layouts.
- `web/src/component/l3/`: product screens such as auth, dashboard, and loading
  states.
- `web/src/lib/`: small non-component browser helpers, including `cx()`.
- `web/src/styles.css`: Tailwind input, source globs, and the `data-theme` dark
  variant. Global `html`/`body` sizing and component CSS are intentionally not
  allowed here; use Tailwind utilities and component tokens instead.
- `web/src/write-index.ts`: writes the generated app shell with hashed asset
  references after Bun builds React.

Generated output such as `web/public/`, `web/src/tailwind.css`, `.carbide/`,
and `web/node_modules/` is ignored.

## Agent Context

Generated apps start without local `AGENTS.md`, `README.md`, or `agents.d/`.
Agent startup guidance is centralized at:

```text
https://carbide.ryangerardwilson.com/for/agents
```

Runtime and framework truth stays in `carbide.toml`, `docker-compose.yml`, and
the `web/`, `api/`, and `db/` source trees. If local `README.md` or
`AGENTS.md` files exist, treat them as app-owned context.
