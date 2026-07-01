# Initial User Experience

The first Sealion experience should feel close to Laravel's default product
loop: install one command, create an app, run one dev command, and land in a
working browser experience with auth already present. The default UI is now a
Bun/React/Tailwind frontend container backed by a Go API container and
Postgres.

## Happy Path

```sh
curl -fsSL https://raw.githubusercontent.com/ryangerardwilson/sealion/main/install.sh | bash
sealion new demo
cd demo
sealion run dev
sealion status
sealion stop dev
```

The installer builds the CLI with Go. The generated app does not require host
Bun, Node, backend Go setup, or Postgres because those run inside containers.

Then open:

```text
the app URL printed by sealion run dev
```

`sealion run dev` prints the working app and API URLs. It prefers port 8080 and
silently selects another local port when 8080 is already in use. To choose one
explicitly:

```sh
SEALION_HTTP_PORT=18080 sealion run dev
```

The frontend listens on port 8080 inside its container. The browser URL is the
host URL printed by the CLI. API calls use the same origin under `/api`.

When Docker Compose supports file watch, `sealion run dev` starts the stack with
quiet Compose output, watch enabled, and live logs streamed below the startup
summary. Edits under `view/web/src/`, `src/`, `model/`, `controller/`, view web
package/config files, `go.mod`, `go.sum`, or `Dockerfile` rebuild and replace
the relevant container.

The CLI presents output as aligned rows and compact tables with TTY-only color
and terminal-only ILoveCandy-style per-container startup and shutdown animation
while Docker Compose builds, starts, waits for, or stops containers. Captured
or piped output remains plain text for tests, scripts, and AI agents. Before
startup, `sealion run dev` prints only the working app and API URLs. Logs begin
only after Compose reports the stack ready. `Ctrl+C` detaches from live log
streaming and leaves the containers running. `sealion follow logs` attaches to
live container logs again. `sealion status` prints the current service table.
`sealion stop dev` stops the local development stack. Frontend, backend,
database, and watch events appear in one timestamped, service-tagged stream and
are mirrored to `.sealion/log/dev.jsonl`. `NO_COLOR` disables ANSI color
without disabling the terminal startup or shutdown animation.

The generated app starts with no seeded users. The first browser visit opens the
account creation flow. Registration creates the first user and session; later
sessions use the login form.

Generated apps keep browser UI in `view/web/src/`. Bun owns the frontend
server and API proxy, Tailwind owns styling, and React owns page flow, forms,
and dashboard rendering. The Go backend owns `/api` routes, auth, sessions,
validation, and Postgres access. The frontend proxies `/api` and `/health` to
the backend so cookies remain same-origin.

The generated app includes:

- a Bun/React/Tailwind frontend container;
- a Go backend/API container;
- a Postgres service container;
- checked-in Docker Compose infrastructure;
- register, login, logout, and dashboard routes;
- model and controller directories for backend code;
- Postgres-backed users and sessions;
- queryable structured dev logs.

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

Runs the generated app through Docker Compose. The frontend, backend, and
database are separate services, matching the runtime topology contract. The CLI
prints the app URL and API URL, then streams logs until `Ctrl+C`. `Ctrl+C`
detaches from the log stream without stopping containers.

### `sealion status`

Prints a table of Compose services, container names, published host ports,
internal container ports, and status.

### `sealion stop dev`

Stops the generated app's Docker Compose dev stack. This is the explicit
container teardown command.

### `sealion logs`

Reads `.sealion/log/dev.jsonl`, the structured log file written by `sealion run
dev`. It supports simple word-based queries such as `sealion logs service
backend`, `sealion logs containing "/api/login"`, and `sealion logs json`.
`sealion follow logs` attaches to live container logs again after detaching.

## Product Principle

The first useful action is not "read docs." The first useful action is a running
app with a database-backed login flow.
