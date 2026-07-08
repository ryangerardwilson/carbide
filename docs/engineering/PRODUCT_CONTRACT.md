# Product Contract

Carbide is a Docker-first monorepo framework for full-stack apps with React,
Go, Postgres, Bun, and Tailwind.

The product exists to stop premature service sprawl. A normal Carbide app
starts as one coherent repository with clear runtime containers:

- `web`: Bun, React, Tailwind, browser assets, and same-origin API proxy.
- `api`: Go HTTP API, auth, sessions, validation, request logs.
- `db`: Postgres data access and checked-in migrations.

## README Boundary

The root README is the framework-agent entrypoint. It should speak directly to
agents and maintainers updating Carbide itself.

The README should answer:

- Am I working on the framework repo or a generated app?
- Where should generated-app agents go? (`https://carbide.ryangerardwilson.com/for/agents`)
- Where is the checked-in source for that public app-agent contract?
- What are Carbide's goals and non-goals?
- Which engineering files own the framework contract?
- Which verification loop proves a framework change?
- Where is the docs website managed and deployed?

The README must not become:

- a phase roadmap,
- a detailed implementation backlog,
- a duplicate of all engineering docs,
- a generated-app operating guide,
- a long-form architecture spec.

App-agent startup guidance belongs in `docs/site/for/agents.md` and is served
publicly at `https://carbide.ryangerardwilson.com/for/agents`. Framework
implementation detail belongs in `docs/engineering/`.

## Product Tone

Use direct product language. Prefer concrete nouns:

- Docker-first monorepo
- generated app
- web/API/db containers
- Postgres-only
- infrastructure as code
- `carbide health`

Avoid vague positioning like "Laravel-inspired" as the lead idea. Laravel can
remain context in docs when useful, but it is not the current product thesis.
