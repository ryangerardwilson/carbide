# Carbide Agent Guide

The canonical agent startup guide is:

```text
https://carbide.ryangerardwilson.com/for/agents
```

Its checked-in source is `docs/site/for/agents.md`. Do not create another
agent runbook, local `agents.d/` tree, or competing install flow. If agent
startup guidance changes, update `docs/site/for/agents.md`, the homepage
summary, and the regression checks in the same change.

## Local Role

This file is only the framework repo entrypoint. It exists to tell agents where
the product truth lives and how to verify repo changes. Framework engineering
contracts live under `docs/engineering/`; public user-facing guidance belongs
in the README and docs site.

## Hard Rules

- Keep Carbide monorepo-first. Containers define runtime boundaries; they do
  not justify default microservice sprawl.
- Keep generated apps Docker-first with `web`, `api`, and `db` as the root
  service directories.
- Keep `/for/agents` as the central point of truth for agent startup.
- Keep README human-first: what Carbide is, what the user gets, how to install,
  and where to read more.
- Do not reintroduce `agents.d` in the framework root, generated scaffold, or
  docs app.
- Do not add root-level generated app files. Generated-app source belongs in
  `scaffold/`; docs-app runtime belongs in `docs/app/`.
- Before changing scaffold output, update contract tests and run the scaffold
  checks.

## Normal Verification

Use the narrowest check that proves the change, then run broader checks before
shipping contract changes:

```sh
cd cli && go test ./...
bash tests/contract/check_repo_contract.sh
PATH=/home/ryan/.local/share/mise/installs/go/1.26.4/bin:$PATH bash tests/scaffold/cli_scaffold.sh
carbide doctor framework
```

For docs app changes:

```sh
cd docs/app/web && bun run typecheck && bun run assets:build
cd docs/app && docker compose build web
```

For scaffold web changes:

```sh
cd scaffold/web && bun install --frozen-lockfile && bun run typecheck && bun run assets:build
```
