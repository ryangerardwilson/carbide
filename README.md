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
  `carbide follow logs`, `carbide logs`, and `carbide doctor` keep the local
  stack understandable.
- **Explicit baselines:** runtime versions are recorded in `carbide.toml`;
  `carbide doctor` rejects floating Docker images, `latest`, and framework-owned
  semver ranges.

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
| `carbide stop dev` | Stop the local development stack. |
| `carbide follow logs` | Stream live container logs. |
| `carbide logs` | Query saved structured dev logs. |
| `carbide doctor` | Check the generated project contract. |
| `carbide doctor env` | Check env and secrets rules without printing secrets. |
| `carbide doctor runtime` | Run the Docker-backed health and auth flow check. |
| `carbide project migrate` | Prepare an AI-assisted framework migration workspace. |
| `carbide deploy check prod` | Classify a deploy target before preview/apply. |
| `carbide deploy preview prod` | Preview a checked-in production deploy target. |
| `carbide deploy apply prod` | Apply a checked-in single-VM production target. |
| `carbide help` | Print the command reference. |
| `carbide version` | Print the installed CLI version and commit. |
| `carbide upgrade` | Upgrade the installed CLI from GitHub. |

`Ctrl+C` in `carbide run dev` detaches from log streaming and leaves containers
running. Use `carbide follow logs` to attach again and `carbide stop dev` to
stop the stack.

## Generated App Layout

```text
my-carbide-app/
|-- AGENTS.md
|-- PROJECT.md
|-- api/
|-- db/
|-- web/
|-- carbide.toml
|-- docker-compose.yml
`-- README.md
```

`web/` is the Bun/React/Tailwind browser container. `api/` is the Go API
container. `db/` owns Postgres schema and database helper code. `PROJECT.md`
owns app-specific product truth. `AGENTS.md` points agents to the central
`/for/agents` guide. The root `carbide.toml` stores the app identity, port
defaults, env contract, runtime baseline, and deploy targets.

At the generated project root, every directory maps to a standalone Docker
service. The root Compose file is the app orchestration file: it describes how
those services run together while each service keeps its own Dockerfile and
dependency files.

## Runtime And Deploy

Carbide records explicit Carbide baselines in `carbide.toml`: Go module
directive, digest-pinned Go/Bun/Debian/Postgres images, exact React and
Tailwind versions, and a runtime contract version. This keeps framework-owned
dependencies boring and reviewable instead of silently floating.

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

## Learn More

- Public docs: <https://carbide.ryangerardwilson.com>
- Engineering docs: `docs/engineering/`
