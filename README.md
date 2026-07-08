# Carbide

If you are changing the Carbide framework repo, start here.

If you are changing a generated Carbide app, stop and use:

<https://carbide.ryangerardwilson.com/for/agents>

That route is the source of truth for app agents.
Its checked-in source is `docs/site/for/agents.md`.

This README is the framework-agent entrypoint and routing layer for Carbide
itself. Narrow framework truth lives in the relevant files under
`docs/engineering/`.
For framework agents, this README is the source of truth for framework agents.

Carbide is a Docker-first monorepo framework for full-stack apps with Bun,
React, Tailwind, Go, and Postgres. A normal Carbide app starts as one repo
with three runtime containers: `web`, `api`, and `db`.

## Goals

- Keep Carbide Docker-first and monorepo-first.
- Keep the generated app shape stable: `web`, `api`, `db`.
- Keep Postgres mandatory and checked in as part of the product contract.
- Keep runtime, env, deploy, and log behavior inspectable by agents.
- Keep the docs app dogfooding the current framework contract.
- Keep framework regressions runner-agnostic and easy to prove locally.

## Non-Goals

- Default microservice sprawl.
- Optional first-party database modes.
- Framework-owned rewrites of app code after scaffold.
- Duplicate root-level framework runbooks.
- GitHub-specific CI or deploy assumptions as product truth.

## Source Of Truth

- Generated app agents:
  `https://carbide.ryangerardwilson.com/for/agents`
- Checked-in app-agent contract source: `docs/site/for/agents.md`
- When the app-agent contract itself needs to change, edit
  `docs/site/for/agents.md`, which is published at
  `https://carbide.ryangerardwilson.com/for/agents`
- Product boundary: `docs/engineering/PRODUCT_CONTRACT.md`
- Repo layout: `docs/engineering/REPO_STRUCTURE.md`
- Generated app contract: `docs/engineering/SCAFFOLD_CONTRACT.md`
- Eternal laws: `docs/engineering/LAWS.md`
- Current starter taste: `docs/engineering/TASTE_GUIDE.md`
- Docs app contract: `docs/engineering/DOCS_APP.md`
- Verification loop: `docs/engineering/REGRESSION_CHECKS.md` and
  `carbide health framework`

Read the smallest relevant engineering file before editing. Do not treat
the public `/for/agents` contract as the framework repo operating guide when
you are editing Carbide itself.

## Task Router

- CLI parsing, output, versioning, upgrade behavior:
  `docs/engineering/CLI_AND_VERSIONING.md`
- Scaffolded app shape, starter files, generated app expectations:
  `docs/engineering/SCAFFOLD_CONTRACT.md`
- Repo-root structure, ownership boundaries, root deletions:
  `docs/engineering/REPO_STRUCTURE.md`
- Docs website runtime, docs deploy behavior, and `/for/agents` publishing:
  `docs/engineering/DOCS_APP.md`
- Product boundary, README scope, framework-vs-app split:
  `docs/engineering/PRODUCT_CONTRACT.md`
- Laws and taste that app agents must follow:
  `docs/engineering/LAWS.md` and `docs/engineering/TASTE_GUIDE.md`

If a framework change alters app-facing commands, laws, taste, deploy
behavior, or agent expectations, update `docs/site/for/agents.md`. That file
is published at `https://carbide.ryangerardwilson.com/for/agents`.

## Verification

Use the narrowest check that proves the change, then run the full framework
loop before shipping contract changes:

```sh
cd cli && go test ./...
bash tests/contract/check_repo_contract.sh
bash tests/scaffold/cli_scaffold.sh
carbide health framework
```

If a required local toolchain is missing from `PATH`, use the narrower
engineering doc that owns that toolchain rather than hardcoding a workstation
path here.

## Docs Website

- Durable public docs content lives in `docs/site/`.
- The docs runtime and deploy target live in `docs/app/`.
- The public app-agent contract lives at
  `https://carbide.ryangerardwilson.com/for/agents`.
- The checked-in source for that contract is `docs/site/for/agents.md`.
- Build docs web assets with:

```sh
cd docs/app/web
bun run typecheck
bun run assets:build
cd ../
docker compose build web
```

- Deploy docs from `docs/app/` only after local checks pass:

```sh
carbide deploy check prod
carbide deploy preview prod
carbide deploy apply prod
```

- `CARBIDE_DOCS_DEPLOY_SSH` must be present in the shell environment or CI
  secrets before deploy.
- After deploy, verify:

```sh
bash tests/smoke/docs_for_agents_http.sh
curl -fsS https://carbide.ryangerardwilson.com/health
```
