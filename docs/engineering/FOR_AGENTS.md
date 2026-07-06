# For Agents

This page is for AI coding agents helping a developer create, inspect, modify,
or deploy a Carbide application.

## How To Proceed

- First, check whether the current directory contains a Carbide app by looking
  for `carbide.toml`, `AGENTS.md`, and the `web/`, `api/`, and `db/`
  directories.
- If it does, do not create a new app. Read `AGENTS.md`, then the smallest
  relevant files under `agents.d/`, and continue with the user's requested
  task.
- If it does not, verify that Docker, Docker Compose, Git, curl, and Go are
  available.
- If the Carbide CLI is missing, install it from the public repository.
- Prefer Carbide defaults unless the user explicitly asks for a different app
  name, port, deploy target, or architecture.

## Prerequisites

Run quick checks:

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

If Docker or Docker Compose is missing, stop and ask the user to install them.
If Go is missing, install Go or ask the user to install it before building the
CLI.

## Install The CLI

```shell
curl -fsSL https://raw.githubusercontent.com/ryangerardwilson/carbide/main/cli/install.sh | bash
carbide version
```

## Create The Application

Use a human app name when the user gives one. Carbide stores the display name
and creates a slug for the directory.

```shell
carbide new "My Carbide App"
cd my-carbide-app
carbide doctor
carbide run dev
```

Use `carbide init` only when the current directory is empty and the user wants
the app created in place.

## Existing Apps

When working inside an existing app:

- Read `AGENTS.md` first.
- Read only the relevant files under `agents.d/`.
- Run `carbide doctor` before changing framework-owned contracts.
- Keep the app as one coherent repository with `web`, `api`, and `db`
  containers.
- Do not introduce service sprawl unless the user explicitly asks for it and
  the deploy target already supports the role split.
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
- Prefer the generated app contract when the user has not specified a
  different approach.
- Keep secrets out of Git. Use the env contract in `carbide.toml` and local
  `.env` files.
- Use `carbide doctor` as the first regression gate and `carbide doctor
  runtime` when runtime behavior matters.
- Keep answers brief and include the exact commands the user can run next.

## Example Outcome

When setup is complete, leave the user with a Carbide app they can open in the
browser, the current local URL, whether the Docker stack is still running, and
which doctor checks passed.
