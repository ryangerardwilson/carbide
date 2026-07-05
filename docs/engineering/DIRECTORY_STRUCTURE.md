# Directory Structure

This page describes the project layout users get after running
`carbide new "My Carbide App"` or `carbide init`. Carbide's own repository
layout is contributor material; the product docs should lead with the generated
app because that is the structure application teams live in.

```text
my-carbide-app/
|-- AGENTS.md
|-- README.md
|-- .env.example
|-- .gitignore
|-- carbide.toml
|-- docker-compose.yml
|-- agents.d/
|   |-- BACKUP_RESTORE.md
|   |-- DEPLOY.md
|   |-- ENVIRONMENT.md
|   `-- TAILWIND_COMPONENTS.md
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
    `-- src/
        |-- component/
        |   |-- l1/
        |   |   |-- Button.jsx
        |   |   |-- Field.jsx
        |   |   |-- Surface.jsx
        |   |   |-- Text.jsx
        |   |   |-- ThemeToggle.jsx
        |   |   |-- index.js
        |   |   `-- tokens.js
        |   |-- l2/
        |   |   |-- AuthForm.jsx
        |   |   |-- Layouts.jsx
        |   |   `-- index.js
        |   `-- l3/
        |       |-- AuthView.jsx
        |       |-- DashboardView.jsx
        |       |-- LoadingView.jsx
        |       `-- index.js
        |-- lib/
        |   `-- cx.js
        |-- main.jsx
        |-- server.jsx
        |-- styles.css
        `-- write-index.mjs
```

## Root Contract

The generated project root is intentionally small:

- `web/`, `api/`, and `db/` map to standalone Docker services.
- `agents.d/` is the only non-container root directory; it stores local agent
  operating context.
- `docker-compose.yml` owns local runtime coordination across services.
- `carbide.toml` owns the app name, slug, default dev port, runtime baseline,
  environment contract, and deploy targets.
- `AGENTS.md` is the generated agent-facing entrypoint.
- `.env.example` documents local development variables without storing real
  secrets.

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

- `web/src/component/l1/`: primitives and Tailwind utility tokens.
- `web/src/component/l2/`: reusable composed patterns such as forms and
  layouts.
- `web/src/component/l3/`: product screens such as auth, dashboard, and loading
  states.
- `web/src/lib/`: small non-component browser helpers, including `cx()`.
- `web/src/styles.css`: Tailwind input, theme variables, and global browser
  rules.
- `web/src/write-index.mjs`: writes the generated app shell with hashed asset
  references after Bun builds React.

Generated output such as `web/public/`, `web/src/tailwind.css`, `.carbide/`,
and `web/node_modules/` is ignored.

## Agent Context

Generated apps include:

- `agents.d/ENVIRONMENT.md` for environment and secrets handling.
- `agents.d/DEPLOY.md` for preview-before-apply deploy behavior.
- `agents.d/BACKUP_RESTORE.md` for Postgres backup and restore context.
- `agents.d/TAILWIND_COMPONENTS.md` for Tailwind component organization.

These files are part of the app contract because Carbide is built for
developer-and-agent workflows. They document how future changes should preserve
the generated app's runtime and frontend boundaries.
