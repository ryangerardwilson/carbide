# Initial User Experience

The first Sealion experience should feel close to Laravel's default product
loop: install one command, create an app, run one dev command, and land in a
working browser experience with auth already present.

## Happy Path

```sh
curl -fsSL https://raw.githubusercontent.com/ryangerardwilson/sealion/main/install.sh | bash
sealion new demo
cd demo
sealion run dev
```

Then open:

```text
the app URL printed by sealion run dev
```

If port 8080 is already in use, `sealion run dev` selects another local port.
To choose one explicitly:

```sh
SEALION_HTTP_PORT=18080 sealion run dev
```

The app listens on port 8080 inside the container. The browser URL is the host
URL printed by the CLI and app logs.

When Docker Compose supports file watch, `sealion run dev` starts the stack with
Compose watch enabled. Edits under `src/`, `views/`, or to `Dockerfile` rebuild
and replace the app container, which matches the current C compile model.

Generated apps keep page markup in `views/*.html`. The starter renderer supports
escaped variables with `{{ name }}` and trusted raw slots with `{!! content !!}`
so developers can edit HTML directly instead of writing C string literals.

The generated app includes:

- a C app container;
- a Postgres service container;
- checked-in Docker Compose infrastructure;
- register, login, logout, and dashboard routes;
- editable HTML views in `views/`;
- Postgres-backed users and sessions;
- a seeded demo account at `admin@sealion.local` with password `password`.

## Commands

### `sealion help`

Prints the command reference.

### `sealion upgrade`

Upgrades the installed CLI when a newer GitHub commit is available.

### `sealion new <project-name>`

Creates a new project directory from the default starter template. It fails if
the target already exists.

### `sealion init`

Initializes the current directory from the default starter template. It fails
unless the current directory is empty.

### `sealion run dev`

Runs the generated app through Docker Compose. The app and database are separate
services, matching the runtime topology contract.

## Product Principle

The first useful action is not "read docs." The first useful action is a running
app with a database-backed login flow.
