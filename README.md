# Carbide

Carbide is a Docker-first monorepo framework for full-stack apps with React,
Go, Postgres, Bun, and Tailwind.

It is built around one product repo with clear runtime containers, not a pile
of premature microservices. A Carbide app starts with a Bun/React/Tailwind web
container, a Go API container, and a mandatory Postgres db container. The
infrastructure, environment contract, logs, and deploy targets live beside the
app code.

## What You Get

- **One app repo:** product code, generated Docker Compose setup, environment
  rules, and deploy targets stay together.
- **Separate runtime boundaries:** `web`, `api`, and `db` run as separate
  containers with separate health checks, logs, and lifecycles.
- **React frontend:** Bun, React, and Tailwind run inside the `web` container,
  not on every developer's host machine.
- **Go API:** auth, sessions, validation, request logging, and application
  routes live in the API container.
- **Postgres-only data:** Carbide treats Postgres as the durable application
  database, not as one adapter among many.
- **Built-in first screen:** registration, login, logout, dashboard, light/dark
  mode, and a small Tailwind component structure are scaffolded on day one.
- **Infrastructure as code:** services, ports, health checks, volumes, env
  values, secrets policy, and deploy targets are checked in.
- **Inspectable dev loop:** `carbide run dev`, `carbide status`,
  `carbide clean dev`, `carbide follow logs`, `carbide logs`, and
  `carbide health` keep the local stack understandable.
- **Starter defaults, not hidden ownership:** the scaffold ships with pinned
  runtime choices, but generated app code belongs to the app immediately.
  Carbide audits laws; Codex makes intentional app edits.

## Quick Start

```sh
curl -fsSL https://raw.githubusercontent.com/ryangerardwilson/carbide/main/cli/install.sh | bash
carbide new "My Carbide App"
cd my-carbide-app
carbide run dev
```

`carbide new "My Carbide App"` creates `my-carbide-app`, stores
`name = "My Carbide App"`, and stores `slug = "my-carbide-app"`.

The installer uses a release binary when one is available and falls back to a
source build when needed. Go is required only for that source-build fallback.
Generated apps run Bun, the Go API build, and Postgres inside containers.
Docker with Docker Compose is required to run generated apps.

## Core Commands

| Command | Use |
| --- | --- |
| `carbide new <project-name>` | Create a new Carbide app directory. |
| `carbide init` | Initialize the current empty directory. |
| `carbide run dev` | Start the local web, API, and Postgres containers. |
| `carbide status` | Show services, containers, ports, and health. |
| `carbide urls` | Show the local app and API URLs. |
| `carbide clean dev` | Normalize local dev state without deleting volumes. |
| `carbide stop dev` | Stop the local development stack. |
| `carbide follow logs` | Stream live container logs. |
| `carbide logs` | Query saved structured dev logs. |
| `carbide health` | Check the generated app laws. |
| `carbide health env` | Check env and secrets rules without printing secrets. |
| `carbide health runtime` | Run the Docker-backed health and auth flow check. |
| `carbide audit` | Start a Codex session to audit the app for compliance with the Carbide contract, and fix any loose ends or missing files. |
| `carbide deploy check prod` | Classify a deploy target before preview/apply. |
| `carbide deploy preview prod` | Preview a checked-in production deploy target. |
| `carbide deploy apply prod` | Apply a checked-in single-VM production target. |
| `carbide help` | Print the command reference. |
| `carbide version` | Print the installed CLI version and commit. |
| `carbide upgrade` | Upgrade the installed CLI from GitHub. |

`Ctrl+C` in `carbide run dev` detaches from log streaming and leaves containers
running. Use `carbide follow logs` to attach again, `carbide stop dev` for
explicit teardown, or `carbide clean dev` when the current session state is
unclear and you want a fresh restart without deleting volumes.

## Generated App Layout

```text
my-carbide-app/
|-- .env.example
|-- api/
|-- db/
|-- web/
|-- carbide.toml
`-- docker-compose.yml
```

`web/` is the Bun/React/Tailwind browser container. `api/` is the Go API
container. `db/` owns Postgres schema and database helper code. The root
`carbide.toml` stores the app identity, port defaults, the env contract,
current starter runtime defaults, and deploy targets. Carbide does not
scaffold `README.md` or `AGENTS.md`; if the app owner creates them later, they
are app-owned prose rather than framework-owned contract files.

At the generated project root, every directory maps to a standalone Docker
service. The root Compose file is the app orchestration file: it describes how
those services run together while each service keeps its own Dockerfile and
dependency files.

## Runtime And Deploy

Carbide records current starter runtime defaults in `carbide.toml`: Go module
directive, digest-pinned Go/Bun/Debian/Postgres images, exact React and
Tailwind versions, and a runtime contract version. Those are starter choices,
not framework ownership over the app after scaffold.

The ownership rule is simple:

- `carbide new` and `carbide init` create the starter.
- After scaffold, the app owns its own code immediately.
- `carbide health` checks eternal laws only.
- `carbide audit` creates a comparison workspace under `.carbide/audit/` and,
  in an interactive terminal with `codex` installed, launches the audit in the
  Codex CLI. Carbide itself does not rewrite app code.
- If app code changes during an audit, those changes are Codex edits made
  intentionally inside the app.

Deploy targets are environments, not just machines. A simple production target
can be one host. A larger target can describe multiple hosts and roles so web,
API, and db can be planned separately.

```sh
carbide deploy check prod
carbide deploy preview prod
carbide deploy apply prod
```

`check` classifies the target as missing, preview-only, invalid, or
apply-supported. `preview` is non-mutating. `apply` is the only command allowed
to change remote infrastructure. Carbide supports `ssh-compose` apply for a
checked-in single-VM target; multi-VM `ssh-compose-environment` targets are
preview-only until clustered orchestration is implemented.

## Docs Website

The public docs website is managed from this repo. Checked-in docs content
lives in `docs/site/`. The docs runtime and deploy target live in `docs/app/`.

Typical docs-management loop:

```sh
cd docs/app/web
bun run typecheck
bun run assets:build
cd ../
docker compose build web
carbide health
```

Deploy from `docs/app` after setting `CARBIDE_DOCS_DEPLOY_SSH`:

```sh
carbide deploy check prod
carbide deploy preview prod
carbide deploy apply prod
```

After deploy, verify the public site and the agent route:

```sh
bash ../../tests/smoke/docs_for_agents_http.sh
curl -fsS https://carbide.ryangerardwilson.com/health
```

## Learn More

- Public docs: <https://carbide.ryangerardwilson.com>
- Engineering docs: `docs/engineering/`
