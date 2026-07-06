# Carbide for Agents

This Markdown is the source of truth for AI agents setting up or working inside
a Carbide app.

If a user asks you to build a new Carbide app, follow this page before making
local choices. If the current directory already contains a Carbide app, do not
create another one.

Fallback raw source if this route is unavailable:

```text
https://raw.githubusercontent.com/ryangerardwilson/carbide/main/docs/site/for/agents.md
```

## Source Precedence

When instructions conflict, use this order:

1. The user's latest explicit instruction.
2. The local generated app `AGENTS.md`.
3. Local `carbide.toml` and `docker-compose.yml`.
4. This `/for/agents` guide.
5. The local generated app `README.md`.
6. Public Carbide documentation pages.
7. Carbide framework repo engineering docs, only when working on Carbide
   itself.

## Identify The Current State

A directory is a Carbide app when it has:

- `carbide.toml`
- `docker-compose.yml`
- `AGENTS.md`
- `PROJECT.md`
- `web/`
- `api/`
- `db/`

When those files and directories exist, work inside the existing app:

```shell
carbide doctor
carbide status
```

Read `AGENTS.md`, `PROJECT.md`, `README.md`, `carbide.toml`, and the files
directly related to the user's task. Do not run `carbide new` or `carbide init`
inside an existing app.

## Prerequisites

Carbide generated apps run Bun, React, Tailwind, Go API builds, and Postgres
inside Docker containers. The host needs Docker, Docker Compose, Git, and curl.
The installer uses a release binary when available and falls back to a source
build, so Go is needed only for that source-build fallback.

Use quick checks when setup is uncertain:

```shell
docker --version
docker compose version
git --version
curl --version
carbide version
```

If Docker or Docker Compose is missing, stop and tell the user Docker is
required. If the installer reports that Go is missing, install Go or ask the
user to install it before using the source-build fallback.

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
carbide clean dev
carbide follow logs
carbide logs
carbide doctor
carbide doctor env
carbide doctor runtime
carbide stop dev
```

`carbide run dev` starts the Docker stack and streams logs after the stack is
ready. `Ctrl+C` detaches from log streaming and leaves containers running. Use
`carbide follow logs` to attach again, `carbide stop dev` for explicit
teardown, or `carbide clean dev` when session state is unclear and you want a
fresh restart without deleting volumes.

Carbide uses command-shaped JSON output, not `--json` flags. Use
`carbide status json`, not `carbide status --json`.

For machine-readable state, use JSON subcommands:

```shell
carbide urls json
carbide status json
carbide doctor json
carbide doctor env json
carbide doctor runtime json
carbide doctor framework json
carbide deploy check prod json
carbide deploy preview prod json
```

Use `carbide help` for the command reference. Use `carbide upgrade` to update
the installed CLI.

## Generated App Contract

Generated apps are Docker-first monorepos:

- `web/` is the Bun, React, Tailwind, and TypeScript browser container.
- `api/` is the Go HTTP/API container.
- `db/` owns Postgres data access and checked-in migrations.
- `docker-compose.yml` owns local service orchestration.
- `carbide.toml` owns app identity, default port, runtime baselines, env
  contract, and deploy targets.
- `AGENTS.md` points agents back to this `/for/agents` guide.
- `PROJECT.md` owns app-specific product truth: domain facts, users, roles,
  business rules, and acceptance criteria.

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

Use these escape hatches instead of expanding global CSS:

- Put reusable class groups in `web/src/component/l1/tokens.ts`.
- Compose variants inside components with TypeScript helpers.
- Keep third-party CSS explicit and product-owned.
- If a product intentionally needs global CSS, create `web/src/product.css`,
  import it explicitly, document the reason in `PROJECT.md`, and update the
  Carbide doctor contract instead of hiding it in `styles.css`.

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
carbide deploy check prod
carbide deploy preview prod
carbide deploy apply prod
```

Carbide supports `ssh-compose` apply for a checked-in single-VM target. New
apps ship with no deploy target, so `carbide deploy apply prod` refuses until a
target exists. `ssh-compose-environment` validates and previews multi-VM
topology, but apply is guarded until clustered orchestration is implemented.

Read `carbide.toml` before assuming topology. Do not mutate infrastructure
without a preview. Use JSON when an agent or CI needs stable state:

```shell
carbide deploy check prod json
carbide deploy preview prod json
```

## Migration And Upgrades

Use `carbide upgrade` for the installed CLI.

When the user asks to move an existing Carbide app toward the latest scaffold
contract, use:

```shell
carbide project migrate
```

`carbide project migrate` creates:

- `.carbide/migration/<timestamp>/latest-scaffold/`: the newest scaffold
  rendered with the current app name and slug.
- `.carbide/migration/<timestamp>/MIGRATION.md`: the migration brief and
  verification loop.

Treat that workspace as a manual framework-upgrade aid, not as a deterministic
code rewrite. The command does not port business logic, merge conflicts, or
rewrite app-owned files for you. Preserve app-specific behavior, data, deploy
targets, secrets, and externally visible behavior while porting
framework-owned files from `latest-scaffold/`.

## Troubleshooting

Use the smallest command that classifies the failure before editing code.

- Doctor or env failures:
  Run `carbide doctor json` or `carbide doctor env json`. Fix the named
  contract first. Missing env values belong in `.env` for local dev or the
  deploy secret layer for remote targets.
- Container start failure or unclear local state:
  Run `carbide clean dev`, then `carbide run dev`. If you need the current
  stack instead of a reset, run `carbide status`, `carbide follow logs`, or
  query `.carbide/log/dev.jsonl` with `carbide logs service api`.
- Deploy, nginx, or sudo failure:
  Run `carbide deploy check prod` and `carbide deploy preview prod` first. If
  Carbide-managed nginx fails, the remote user needs non-interactive sudo for
  nginx install/reload, or the target should set `nginx = false` and use
  user-managed ingress.
- Version drift or framework mismatch:
  Run `carbide version`. For framework work, run `carbide doctor framework`.
  For generated apps, upgrade the CLI, run `carbide project migrate`, then
  port framework-owned files manually against `latest-scaffold/`.

## Verification

Use the smallest check that proves the change, then widen before finishing.

Fast first-run verification:

```shell
carbide doctor
carbide status
```

Full runtime verification when Docker is available or the task changed
container, API, auth, or cache behavior:

```shell
carbide doctor runtime
```

Common app checks:

```shell
cd web && bun run typecheck && bun run assets:build
cd ../api && go test ./...
```

If the change affects containers, also run a Docker build or the relevant
runtime flow. Report which checks passed and which checks could not be run.

## Recovery

- If `carbide run dev` was interrupted, assume containers are still running.
  Run `carbide status`, then `carbide follow logs`, `carbide stop dev`, or
  `carbide clean dev`.
- If ports are unclear, run `carbide urls` or `carbide urls json`.
- If deploy is unclear, run `carbide deploy check prod`.
- If a doctor check fails, fix the named contract before adding new behavior.
- If a task needs product facts that are not in source, update `PROJECT.md`
  after the user confirms them.

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
