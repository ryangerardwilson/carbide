# Carbide for Agents

This page is for AI coding agents that are helping a developer start or work
inside a Carbide application.

## How to Proceed

- First, check whether the current directory already contains a Carbide app.
- A Carbide app has `carbide.toml`, `AGENTS.md`, and `web/`, `api/`, and
  `db/` directories.
- If it does, skip installation and move straight to the user's requested
  task.
- If it does not, verify that Docker, Docker Compose, Git, curl, Go, and the
  Carbide CLI are available.
- If the Carbide CLI is missing, install it from GitHub.

## Prerequisites

Run quick version checks:

```shell
docker --version
docker compose version
git --version
curl --version
go version
carbide version
```

Host Bun, Node, React, Go API setup, and Postgres setup are not prerequisites
for a generated Carbide app. They run inside Docker containers.

If Docker or Docker Compose is missing, tell the user that Docker is required
before continuing.

If Go is missing, install Go or ask the user to install it before building the
CLI.

## Install the CLI

```shell
curl -fsSL https://raw.githubusercontent.com/ryangerardwilson/carbide/main/cli/install.sh | bash
carbide version
```

## Create the Application

Create the new application with Carbide defaults.

```shell
carbide new "My Carbide App"
cd my-carbide-app
carbide doctor
carbide run dev
```

Use `carbide init` only when the current directory is empty and the user wants
the app created in place.

## Existing Apps

When the current directory is already a Carbide app:

- Read `AGENTS.md` first.
- Read only the relevant files under `agents.d/`.
- Run `carbide doctor` before changing framework-owned contracts.
- Keep the app as one coherent repository with `web`, `api`, and `db`
  containers.
- Do not introduce service sprawl unless the user explicitly asks for it and
  the deploy target already supports that topology.
- Use `carbide project migrate` when the app needs to be moved toward a newer
  scaffold contract.

## Development Loop

```shell
carbide run dev
carbide status
carbide follow logs
carbide logs containing "/api/login"
carbide doctor
carbide doctor runtime
carbide stop dev
```

`carbide run dev` starts Docker containers and streams logs after the stack is
ready. `Ctrl+C` detaches from logs and leaves containers running. Use
`carbide stop dev` to stop them.

## Deployment

Inspect deployment targets before mutating a server:

```shell
carbide deploy preview prod
carbide deploy apply prod
```

A deploy target is an environment. It may be a single VM or a set of hosts and
roles. Read `carbide.toml` before assuming topology.

## Guidance

- Ask only for decisions that materially affect the app.
- Prefer Carbide defaults when the user has not specified a preference.
- Keep secrets out of Git. Use the env contract in `carbide.toml` and local
  `.env` files.
- Use `carbide doctor` as the first regression gate.
- Use `carbide doctor runtime` when runtime behavior matters.
- Keep answers brief and include exact commands the user can run next.

## Example Outcome

When setup is complete, leave the user with a Carbide app they can open in the
browser, the current local URL, whether the Docker stack is still running, and
which doctor checks passed.
