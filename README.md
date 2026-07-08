# Carbide

If you are changing a generated Carbide app, stop and use:

<https://carbide.ryangerardwilson.com/for/agents>

That route is the source of truth for app agents.
Its checked-in source is `docs/app/web/site/for/agents/index.md`.

If you are changing the Carbide framework repo itself, this `README.md` is the
framework-agent entrypoint and routing layer.
For framework agents, this README is the source of truth for framework agents.

Carbide is a Docker-first monorepo framework for full-stack apps with Bun,
React, Tailwind, Go, and Postgres. A normal Carbide app starts as one repo
with three runtime containers: `web`, `api`, and `db`.

There is no separate internal docs tree under `docs/engineering/`.
Framework truth lives in this README, the checked-in public docs pages under
`docs/app/web/site/`, the scaffold under `scaffold/`, the CLI source under
`cli/`, and the executable contracts under `tests/`.

## Goals

- Keep Carbide Docker-first and monorepo-first.
- Keep the generated app shape stable: `web`, `api`, `db`.
- Keep Postgres mandatory.
- Keep app-facing laws explicit and executable through `carbide health`.
- Keep the docs app dogfooding the current framework contract.
- Keep framework regressions runner-agnostic and provable locally.

## Non-Goals

- Default microservice sprawl.
- Optional first-party database modes.
- Framework-owned rewrites of app code after scaffold.
- Duplicate framework runbooks or shadow contract trees.
- GitHub-specific CI or deploy assumptions as product truth.

## Source Of Truth

- Generated app agents:
  `https://carbide.ryangerardwilson.com/for/agents`
- Checked-in app-agent contract source:
  `docs/app/web/site/for/agents/index.md`
- Framework agents:
  this `README.md`
- Human docs page sources:
  `docs/app/web/site/index.html`,
  `docs/app/web/site/create-your-first-app.html`,
  `docs/app/web/site/frontend-starter-contract.html`,
  `docs/app/web/site/repo-structure.html`,
  `docs/app/web/site/deployment.html`,
  `docs/app/web/site/ci-cd-regression-tests.html`,
  `docs/app/web/site/version-policy.html`
- Generated starter contract:
  `scaffold/`
- CLI/runtime implementation:
  `cli/`
- Executable framework contract:
  `tests/` and `carbide health framework`

If a framework change alters app-facing commands, laws, taste, docs, or deploy
behavior, update `docs/app/web/site/for/agents/index.md`, the relevant public
docs page under `docs/app/web/site/`, the scaffold, and the tests together.

## Task Router

- CLI parsing, health, audit, deploy, logs, upgrade:
  `cli/internal/cli/`
- Installer and installed CLI wrapper:
  `cli/install.sh`, `cli/bin/carbide`
- Generated app shape and starter files:
  `scaffold/`
- Public docs content:
  `docs/app/web/site/`
- Docs runtime and deploy script:
  `docs/app/`
- Framework regressions:
  `tests/`

## Current App Laws

When you change generated-app laws, update `/for/agents`, `carbide health`,
and the related tests in the same change.

1. One app repo.
2. Root runtime directories are `web/`, `api/`, and `db/`.
3. `carbide.toml` and `docker-compose.yml` stay checked in.
4. Browser API traffic stays same-origin through `web -> /api -> api`.
5. Postgres is required.
6. Deploy targets point to checked-in scripts inside the app repo.
7. Secrets are never printed.
8. Repo-owned source, config, and test files stay at 1000 lines or fewer.

## Current Starter Taste

When you change starter taste, update `/for/agents`, the public docs pages,
the scaffold, and the related tests together.

- `web/` is a Bun, React, Tailwind, and TypeScript container.
- `api/` is a Go API container.
- `db/` owns Postgres access and checked-in migrations.
- Frontend components are organized as `l1`, `l2`, and `l3`.
- The starter ships register, login, logout, dashboard, and theme behavior.
- Product-owned palettes are allowed. The docs app is intentionally black and yellow
  across light and dark modes, and audits should preserve that.

## Verification

Use the narrowest check that proves the change, then run the full framework
loop before shipping contract changes:

```sh
cd cli && go test ./...
bash tests/contract/check_line_limits.sh
bash tests/contract/check_repo_contract.sh
bash tests/scaffold/cli_scaffold.sh
carbide health framework
```

If a required local toolchain is missing from `PATH`, prefer containerized or
checked-in fallback paths over hardcoded workstation paths.

## Docs Website

- Durable public docs content lives in `docs/app/web/site/`.
- The docs runtime and app-owned deploy script live in `docs/app/`.
- The docs app does not carry its own `AGENTS.md` or `README.md`.
- The public app-agent contract lives at
  `https://carbide.ryangerardwilson.com/for/agents`.
- The checked-in source for that contract is
  `docs/app/web/site/for/agents/index.md`.
- The docs deploy script is `docs/app/deploy/prod.sh`.
- `CARBIDE_DOCS_DEPLOY_SSH` and `CARBIDE_DOCS_POSTGRES_PASSWORD` must be
  present before deploy.

Build docs web assets with:

```sh
cd docs/app/web
bun run typecheck
bun run assets:build
cd ../
docker compose build web
```

Deploy docs from `docs/app/` only after local checks pass:

```sh
carbide deploy prod
```

After deploy, verify:

```sh
bash tests/smoke/docs_for_agents_http.sh
curl -fsS https://carbide.ryangerardwilson.com/health
```
