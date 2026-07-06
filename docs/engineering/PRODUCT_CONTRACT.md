# Product Contract

Carbide is a Docker-first monorepo framework for full-stack apps with React,
Go, Postgres, Bun, and Tailwind.

The product exists to stop premature service sprawl. A normal Carbide app
starts as one coherent repository with clear runtime containers:

- `web`: Bun, React, Tailwind, browser assets, and same-origin API proxy.
- `api`: Go HTTP API, auth, sessions, validation, request logs.
- `db`: Postgres data access and checked-in migrations.

## README Boundary

The README must answer human questions:

- What is Carbide?
- Why would I use it?
- What do I get after `carbide new`?
- What do I need installed?
- What commands do I run?
- Where are the docs?

The README must not become:

- a phase roadmap,
- a detailed implementation backlog,
- a regression checklist,
- an agent instruction file,
- a long-form architecture spec.

Move that material into `docs/engineering/`. Agent startup guidance belongs in
`docs/site/for/agents.md` and is served at `/for/agents`.

## Product Tone

Use direct product language. Prefer concrete nouns:

- Docker-first monorepo
- generated app
- web/API/db containers
- Postgres-only
- infrastructure as code
- `carbide doctor`

Avoid vague positioning like "Laravel-inspired" as the lead idea. Laravel can
remain context in docs when useful, but it is not the current product thesis.
