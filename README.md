# Carbide

Carbide is an experimental, Laravel-inspired full-stack framework with a React
web container, a Go API container, and a mandatory Postgres database.

The goal is not to copy Laravel line by line. The goal is to find the smallest
set of conventions, tools, and runtime guarantees that make building web apps in
containers feel coherent, productive, and safe enough to be practical.

## Product Bet

The product bet is Docker-first convention over host setup: React owns the
browser, Go owns the application API, and Postgres owns durable relational
state. Carbide should make that full-stack default feel boring, inspectable, and
fast to start:

- one Bun/React/Tailwind `web` container
- one Go API container
- one mandatory Postgres `db` container
- one project layout
- one API request lifecycle
- one checked-in schema and migration path
- one checked-in infrastructure contract
- one environment and secrets contract
- one Go CLI entry point
- one conservative set of auth and session defaults

Carbide should make the hard parts visible instead of hiding them behind magic.

## Core Principles

- **Container-first:** every app runs through generated containers, not host
  Bun, host API toolchains, or hidden local services.
- **Go CLI:** `carbide` is a compiled Go CLI. It owns scaffolding, upgrades,
  local port selection, structured terminal output, queryable dev logs, and the
  Docker Compose development lifecycle.
- **React default frontend:** the browser UI lives in the `web` container.
  Bun, React, and Tailwind are required inside that container, not on the
  developer's host machine.
- **Go API:** auth, sessions, validation, API routes, and business logic live
  in the API container.
- **Postgres-only:** Carbide targets Postgres as the mandatory database, not as
  one interchangeable adapter among many.
- **Separate runtime boundaries:** web, API, and db containers
  are separate services with separate lifecycles, health checks, logs, and
  storage.
- **Infrastructure as code:** every supported runtime dependency, service
  boundary, volume, network, secret contract, environment variable, health
  check, and deploy target must be described in checked-in code.
- **Preview before apply:** deploy tooling must show infrastructure changes
  before it is allowed to mutate real resources.
- **Explicit ownership:** requests, responses, sessions, and database handles
  must have clear lifetimes.
- **Convention over configuration:** defaults should cover normal apps without
  requiring boilerplate.
- **Safe by default:** shipped routing, sessions, cookies, validation, SQL
  access, and future upload/CSRF surfaces should have conservative defaults.
- **Inspectable runtime:** generated files, migrations, logs, and app state
  should be easy to inspect and reproduce.
- **Small ecosystem surface:** add extension points only after the core app loop
  is stable.

## Non-Goals

- Running generated apps directly on host Bun, Go, or Postgres before the
  container contract is stable.
- Full Laravel API compatibility.
- Requiring host-installed Bun, Node, or npm.
- Rebuilding React, Bun, or Tailwind from scratch.
- A general-purpose language package manager.
- ORM magic that hides SQL, migrations, or operational behavior.
- Supporting multiple databases, web servers, or production deployment targets
  in the first versions.

## Runtime Topology

The default Carbide app runs as three containers:

1. the `web` container, which owns Bun, React, Tailwind, browser routing,
   the API proxy, and the public host port;
2. the API container, which owns Go API routes, auth, sessions, application
   code, request logs, and startup checks;
3. the Postgres `db` container, which owns durable relational state through
   a mounted volume or managed persistent storage, plus checked-in schema
   state under `db/migration`.

The browser talks to the web service on one origin. The web service proxies `/api` and
`/health` to the API service over the private Compose network, which keeps
cookies same-origin and avoids CORS as the default development problem. The API
depends on Postgres readiness, but each service remains independently
restartable, inspectable, and replaceable.

## Infrastructure As Code Contract

Carbide apps must be reproducible from the repository. Runtime behavior should
not depend on manual console setup, undocumented shell history, or hidden
machine state.

The first supported infrastructure target is a generated Docker Compose setup
for local development. Production targets come later, one at a time, after the
local app and Postgres contract is stable.

At minimum, each generated app keeps these contracts in version control:

- container definitions for the web, API, db, and required
  services;
- service networking, health checks, restart policy, and readiness rules;
- Postgres image version, volume, backup, restore, and migration policy;
- environment variable contract with required, optional, and secret values;
- generated local Compose manifests first, then deployment manifests for each
  supported production target as those targets become official;
- an environment contract that marks required, optional, secret,
  browser-exposed, and framework-owned values;
- project and environment contract version markers for future infrastructure
  changes.

The Carbide CLI generates and validates these files instead of asking
developers to maintain ad hoc infrastructure by hand. Infrastructure is part of
the application source, and changes to it must be reviewable, diffable, and
recoverable.

Local development secrets and deployment secrets stay separate. The generated
Docker Compose stack uses obvious local-only defaults so a new app boots with
`carbide run dev`. Real deployment secrets belong to the deployment or IaC
layer; Carbide does not add a separate secrets container by default.

Deploy commands follow the preview/apply contract:

```sh
carbide deploy preview dev
carbide deploy apply dev
```

`preview` is non-mutating and shows the planned change set. `apply` is the only
path allowed to mutate infrastructure. Generated starter apps still refuse
`apply` until a checked-in deploy target exists. The first concrete target lives
in `docs/app` as `de-sci`, an `ssh-compose` deployment for the documentation
app.

Deploy targets are modeled as environments, not just machines. The simplest
environment can still be one host with `type = "ssh-compose"`. Larger
environments use checked-in hosts and roles with
`type = "ssh-compose-environment"` so `web`, `api`, and `db` can be planned
across different servers. Carbide previews that topology today and keeps
clustered `apply` guarded until migration order, health gates, load balancer
updates, and rollback semantics are explicit.

## Documentation And Automation

The public documentation site is served by the Carbide docs app on `de-sci` at:

```text
https://carbide.ryangerardwilson.com
```

CI currently runs shell syntax checks, Go CLI tests, repository contract checks,
CLI scaffold checks, and a generated-app Docker smoke flow. The broader
regression plan lives in `docs/engineering/CI_CD_REGRESSION_TESTS.md`.
The current repo layout lives in `docs/engineering/DIRECTORY_STRUCTURE.md`.
The frontend starter contract lives in
`docs/engineering/FRONTEND_STARTER_CONTRACT.md`.
The runtime baseline and upgrade policy lives in
`docs/engineering/VERSION_POLICY.md`.

## Install And Start

```sh
curl -fsSL https://raw.githubusercontent.com/ryangerardwilson/carbide/main/cli/install.sh | bash
carbide new demo
cd demo
carbide run dev
carbide status
carbide stop dev
```

The installer currently builds the `carbide` CLI with Go, so Go must be
available on the host machine. Generated apps still run Bun, the Go API
build, and Postgres inside containers. Docker with Docker Compose is required
to run generated apps.

`carbide new <project-name>` creates a new project directory. Human names are
accepted: `carbide new "My Carbide App"` creates `my-carbide-app`, stores
`name = "My Carbide App"`, and stores `slug = "my-carbide-app"`.
`carbide init` initializes the current directory only when it is empty.
`carbide run dev` starts the generated web, API, and Postgres
containers with register, login, logout, and dashboard already wired. It
prints the working app and API URLs, preferring `http://localhost:8080` and
silently selecting another local port when 8080 is already in use. Set
`CARBIDE_HTTP_PORT=<port>` to choose the host port explicitly.

`Ctrl+C` in `carbide run dev` detaches from live log streaming and leaves the
containers running. `carbide follow logs` attaches to live container logs again.
`carbide status` prints a table of Compose services, container names, published
host ports, internal container ports, and status. `carbide stop dev` stops the
local development stack. `carbide help` prints the command reference.
`carbide upgrade` upgrades the installed CLI when a newer GitHub commit is
available. `carbide logs` reads the structured dev log file written by
`carbide run dev`; examples include `carbide logs service api` and
`carbide logs containing "/api/login" json`.
`carbide doctor` runs the fast project contract check: root shape,
`carbide.toml`, Compose services, env/secrets rules, web/API/db contracts,
agent docs, and legacy-regression markers.
`carbide doctor env` validates the env contract in `carbide.toml`, `.env`,
local defaults, and secret/browser exposure rules without printing secret
values. `carbide doctor runtime` runs the heavier Docker-backed health and auth
flow check. `carbide doctor framework` runs framework source regressions from a
Carbide source checkout. `carbide deploy preview <target>` prints the
non-mutating deploy plan, while `carbide deploy apply <target>` runs only for
checked-in deploy targets such as the docs app `de-sci` target; otherwise it
remains guarded. Multi-server environment targets are previewable and
validated, but `apply` is intentionally guarded until clustered orchestration is
implemented.

Runtime versions are treated as explicit Carbide baselines, not hidden moving
dependencies. New projects record the current baseline in `carbide.toml` under
`[runtime]`: supported Go module directive, digest-pinned Go/Bun/Debian/Postgres
images, exact React and Tailwind package versions, and a runtime contract
version. `carbide doctor` fails on floating Docker images, `latest`, semver
ranges for framework-owned web packages, or unsupported Go directives. The
scheduled dependency audit reports newer stable releases and changed image
digests without editing project files.

When Docker Compose supports file watch, `carbide run dev` starts the stack with
quiet Compose output, watch enabled, and live logs streamed below the startup
summary. Edits under `web/src/`, web package/config files,
`api/`, `db/`, or `api/Dockerfile` rebuild and replace the relevant
container.

CLI output is rendered through a small Go output layer: headings, aligned
labels, compact tables, TTY-only color, full-width terminal-only
ILoveCandy-style per-container startup and shutdown animation, timestamped log
rows, and plain text when piped or captured by scripts. `carbide run dev`
prints only the working app/API URLs before the startup animation and log
stream. Logs begin only after Compose reports the stack ready. `NO_COLOR`
disables ANSI color without disabling the terminal startup or shutdown
animation. Every streamed web, API, db, and watch event is also
written as JSONL to
`.carbide/log/dev.jsonl` so humans, scripts, and AI agents can inspect or query
the whole local system from one command.

Generated apps use explicit runtime boundaries. `web/` owns the Bun server,
Tailwind build, content-hashed browser assets, same-origin `/api` calls, and a
small React starter under `component/l1`, `component/l2`, and `component/l3`.
`api/` owns the Go HTTP server, auth, sessions, routing, and JSON responses.
`db/` owns Postgres-backed data access and checked-in migration state.
`carbide.toml` owns the project metadata, default dev port, environment
contract, and deploy guardrails. `AGENTS.md` is the generated agent-facing
entrypoint. `agents.d/` stores the specific operating notes for env, deploy,
backup, restore, and Tailwind component organization.

At the generated project root, every directory except `agents.d/` maps to a
standalone Docker service: `web/`, `api/`, and `db/`. Shared runtime
coordination lives in `docker-compose.yml`.

Carbide scaffolds L1/L2/L3 React components as a starter convention for
Tailwind class ownership: primitives, reusable patterns, and product screens.
That convention teaches the generated app shape; it is not a package ecosystem
or a permanent design-system mandate. Teams can reorganize the frontend with
their own taste, design system, or AI workflow after the app is generated.

The generated `DashboardLayout` is a conventional app shell: persistent left
sidebar, section navigation, account/logout footer, and a main work area.

## Roadmap

### Current Baseline

- Source install from `cli/install.sh`.
- Compiled Go CLI with `new`, `init`, `run dev`, `status`, `stop dev`,
  `follow logs`, `logs`, `doctor`, `doctor env`, `doctor runtime`,
  `doctor framework`, `deploy preview`, and guarded `deploy apply`.
- Generated `web`, `api`, and `db` services with Docker Compose watch.
- Bun/React/Tailwind browser app with register, login, logout, dashboard, and
  left-sidebar app shell.
- Go API backed by Postgres users and sessions.
- Environment/secrets contract in `carbide.toml`.
- Generated `AGENTS.md` and `agents.d` operating notes.
- Queryable structured dev logs in `.carbide/log/dev.jsonl`.
- Fast project doctor and Docker-backed runtime doctor.
- CI coverage for shell syntax, Go CLI tests, repo contract, scaffold checks,
  and generated Docker smoke flow.

### Phase 1: HTTP Core

- Harden routing for common HTTP methods beyond the starter routes.
- Add request parsing helpers for headers, query params, path params, and
  forms.
- Add response helpers for text, JSON, redirects, files, and errors.
- Add middleware chaining with predictable ownership rules.
- Add structured error pages for development and safe production errors.

### Phase 2: Application Kernel

- Harden the generated `web/`, `api/`, and `db/` directory
  contract.
- Expand configuration loading from environment and checked-in defaults.
- Harden environment contract validation and protected framework-owned keys.
- Add service registration without hidden reflection.
- Add logging with request IDs.
- Add graceful shutdown and worker lifecycle hooks.

### Phase 3: Frontend And Assets

- Keep the Bun/React web container as the public local-development
  entrypoint.
- Proxy `/api` and `/health` to the Go API service to preserve same-origin cookies.
- Ship a small L1/L2/L3 React starter for auth, dashboard, and the app shell
  without turning Carbide into a frontend package ecosystem.
- Make Tailwind the mandatory generated styling path.
- Serve the React shell with content-hashed JS and CSS assets by default.

### Phase 4: Database Layer

- Keep Postgres as the required database.
- Harden connection pooling.
- Add migrations with up/down support.
- Add a query builder with parameter binding by default.
- Add schema inspection helpers for Postgres-specific capabilities.
- Explore a constrained database access layer without pretending every
  Eloquent pattern maps cleanly into a containerized Go API service.

### Phase 5: Web App Essentials

- Harden the default generated auth experience: register, login, logout, and
  dashboard.
- Add signed cookies and encrypted session storage.
- Add CSRF protection.
- Add validation primitives.
- Replace the starter password hash with a production-grade password hashing
  contract.
- Add file upload handling with size and type controls.

### Phase 6: Background Work

- Add Postgres-backed queues.
- Add scheduled jobs.
- Add mail driver contracts.
- Add cache contracts.
- Add retries, dead-letter behavior, and job inspection commands.

### Phase 7: Developer Experience

- Harden the `carbide` Go CLI.
- Harden project scaffolding.
- Add migration generation.
- Add infrastructure generation, validation, and diff commands.
- Add test helpers for HTTP requests and database state.
- Add containerized watch/rebuild workflow.
- Add debug tooling for request lifecycle and connection leaks.

### Phase 8: Production Contract

- Define the official production image.
- Define the first production infrastructure-as-code target after local Compose
  is stable.
- Add health checks and readiness checks.
- Extend the existing structured dev logs into the production container
  contract.
- Harden single-VM deploy apply beyond the docs app target.
- Implement clustered apply for previewed multi-VM environment targets.
- Add backup, restore, and migration rollback guidance.

### Phase 9: Ecosystem

- Stabilize extension points.
- Add first-party packages only where the core framework has repeated evidence.
- Document compatibility rules.
- Publish upgrade guides between framework versions.

## Current Milestone

The current milestone is a generated, containerized app that can:

1. boot with one command,
2. serve a Bun/React/Tailwind browser app,
3. proxy same-origin `/api` calls to the Go API container,
4. register the first user,
5. log in, log out, and open the dashboard,
6. connect to the required Postgres container,
7. write queryable structured dev logs,
8. report container status,
9. detach from logs without stopping containers,
10. stop the dev stack explicitly.

That milestone proves the core local loop before Carbide adds production deploy
targets, queues, a richer migration runner, or higher-level database features.
