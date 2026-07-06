# Carbide for Agents

This Markdown is the source of truth for AI agents setting up or working inside
a Carbide app.

If a user asks you to build a new Carbide app, follow this page before making
local choices. If the current directory already contains a Carbide app, do not
create another one.

## Identify The Current State

A directory is a Carbide app when it has:

- `carbide.toml`
- `docker-compose.yml`
- `AGENTS.md`
- `web/`
- `api/`
- `db/`

When those files and directories exist, work inside the existing app:

```shell
carbide doctor
carbide status
```

Read `AGENTS.md`, `README.md`, `carbide.toml`, and the files directly related
to the user's task. Do not run `carbide new` or `carbide init` inside an
existing app.

## Prerequisites

Carbide generated apps run Bun, React, Tailwind, Go API builds, and Postgres
inside Docker containers. The host needs Docker, Docker Compose, Git, curl, and
Go for the CLI installer.

Use quick checks when setup is uncertain:

```shell
docker --version
docker compose version
git --version
curl --version
go version
carbide version
```

If Docker or Docker Compose is missing, stop and tell the user Docker is
required. If Go is missing, install Go or ask the user to install it before
building the CLI.

## Create A New App

For a new app, run:

```shell
curl -fsSL https://raw.githubusercontent.com/ryangerardwilson/carbide/main/cli/install.sh | bash
carbide new demo
cd demo
carbide run dev
carbide doctor
carbide status
```

Human app names are accepted:

```shell
carbide new "My Carbide App"
cd my-carbide-app
```

That creates `my-carbide-app` while storing the display name as
`My Carbide App`.

Use `carbide init` only when the current directory is empty and the user wants
the app created in place.

## Development Loop

Use these commands from the generated app root:

```shell
carbide run dev
carbide status
carbide follow logs
carbide logs
carbide doctor
carbide doctor env
carbide doctor runtime
carbide stop dev
```

`carbide run dev` starts the Docker stack and streams logs after the stack is
ready. `Ctrl+C` detaches from log streaming and leaves containers running. Use
`carbide follow logs` to attach again and `carbide stop dev` to stop the stack.

Use `carbide help` for the command reference. Use `carbide upgrade` to update
the installed CLI when a newer GitHub commit is available.

## Generated App Contract

Generated apps are Docker-first monorepos:

- `web/` is the Bun, React, Tailwind, and TypeScript browser container.
- `api/` is the Go HTTP/API container.
- `db/` owns Postgres data access and checked-in migrations.
- `docker-compose.yml` owns local service orchestration.
- `carbide.toml` owns app identity, default port, runtime baselines, env
  contract, and deploy targets.
- `AGENTS.md` points agents back to this `/for/agents` guide.

Generated apps do not include `agents.d/`. Do not create a local agent runbook
that competes with this source of truth.

The first browser visit should create the first user through the registration
flow. Do not add seeded demo credentials.

## Frontend Contract

The generated web app uses Tailwind as the required styling path.

- Keep `web/src/styles.css` small: Tailwind import, `@source` globs, and the
  `data-theme` dark variant.
- Do not put global `html` or `body` sizing, component CSS, scrollbar CSS, or a
  generated color-variable palette in `web/src/styles.css`.
- Keep repeated visual choices in Tailwind utility tokens and components.
- Preserve visible focus states and semantic form labels.
- Keep the built-in light/dark/system theme behavior browser-local.

The starter uses `web/src/component/l1`, `l2`, and `l3` as Tailwind component
organization. Work with that structure unless the user explicitly asks to
replace it.

## Environment And Secrets

Local development uses Docker Compose defaults and `.env.example` for optional
local overrides. `.env` is ignored.

`carbide.toml` is the checked-in environment contract. Run:

```shell
carbide doctor
carbide doctor env
```

Do not print secret values in logs, docs, CLI output, errors, or chat
responses. Do not add a separate secrets container unless the user explicitly
asks for a new architecture and accepts the operational cost.

## Deployment

Deployment targets live in `carbide.toml`. Preview before applying:

```shell
carbide deploy preview prod
carbide deploy apply prod
```

A target is an environment. It may be a single VM or a set of hosts and roles.
Read `carbide.toml` before assuming topology. Do not mutate infrastructure
without a preview.

## Migration And Upgrades

Use `carbide upgrade` for the installed CLI.

When the user asks to move an existing Carbide app toward the latest scaffold
contract, use:

```shell
carbide project migrate
```

Treat the generated migration workspace as an AI-assisted comparison target,
not as a deterministic code rewrite. Preserve app-specific behavior, data,
deploy targets, secrets, and public domain behavior.

## Verification

Use the smallest check that proves the change, then widen before finishing.

Common checks:

```shell
carbide doctor
carbide doctor runtime
cd web && bun run typecheck && bun run assets:build
cd ../api && go test ./...
```

If the change affects containers, also run a Docker build or the relevant
runtime flow. Report which checks passed and which checks could not be run.

## Agent Behavior

- Prefer Carbide defaults when the user has not specified a choice.
- Ask only for decisions that materially affect the app.
- Keep the app as one coherent monorepo with `web`, `api`, and `db` services.
- Do not introduce extra services, packages, or hosted dependencies without a
  clear user request.
- Keep answers brief and include the exact next command when the user needs to
  act.
- If setup completes, tell the user the local URL, whether containers are still
  running, and which checks passed.
