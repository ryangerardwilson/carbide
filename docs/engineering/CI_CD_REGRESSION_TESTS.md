# CI/CD Regression Test Plan

Carbide's regression plan is runner-agnostic. Run it locally, wire it into any
CI system you trust, or both. The checked-in contract lives in this repo rather
than in GitHub-specific workflow files.

## Required Gates

### Change Gate

Runs before merging or publishing any framework change:

- repository contract checks;
- non-mutating dependency and image drift audit;
- Go CLI unit tests;
- shell syntax checks for repo-owned test and launcher scripts;
- documentation site contract checks;
- generated Docker stack smoke test with registration-first, Postgres-backed
  JSON auth;
- `carbide health` fast app-law checks;
- future API unit, integration, and compatibility checks.

### Publish Gate

Runs before updating shared docs, installers, or branch tips that other people
consume:

- the change gate;
- documentation deployment through the Carbide docs app;
- documentation routes are served as extensionless canonical URLs, with legacy
  `.html` paths redirected;
- future release candidate smoke checks.

### Release Gate

Runs before a deliberate framework release or installer refresh:

- supported Go version matrix;
- race and compatibility matrix;
- Postgres integration matrix;
- generated project smoke tests;
- migration up/down tests;
- container build and boot tests;
- documentation deploy preview or published docs verification.

## Regression Suites

### Repository Contract

Purpose: make the repo shape itself hard to accidentally break.

Initial checks:

- required directories exist;
- README keeps the core product contracts: Bun/React/Tailwind web
  container, Go API container, Postgres-only database, infrastructure as
  code, local Compose first, and Postgres-backed queues;
- install script, CLI, and default scaffold files exist;
- the Go CLI builds and its deterministic helpers and output renderer pass unit
  tests;
- documentation site files exist;
- documentation deployment is owned by the Carbide docs app, not GitHub Pages;
- docs-site internal links use extensionless routes, not `.html` hrefs;
- the repo remains CI-runner-agnostic with no checked-in GitHub workflow
  dependency.

### API Build And Behavior

Purpose: catch broken API builds, incompatible API behavior, and build
drift.

Future checks:

- build generated API code with the pinned Go version;
- run API unit tests with strict failure behavior;
- verify public API routes, cookies, and JSON response shapes;
- fail on accidental generated API contract changes outside explicit release
  work.

### Unit Tests

Purpose: keep core API behavior deterministic.

Future checks:

- router matching;
- request parsing;
- response generation;
- middleware ordering;
- configuration parsing;
- logging shape.

### Compatibility And Race Tests

Purpose: make concurrency and generated app compatibility failures visible
early.

Future checks:

- Go race detector for API packages that can run without containers;
- generated app compatibility tests across supported Go versions;
- request lifecycle tests under concurrent auth and session traffic;
- failure artifacts for reproducible debugging.

### Postgres Integration

Purpose: enforce the mandatory database contract.

Future checks:

- API connects only after Postgres readiness;
- connection pool opens and closes cleanly;
- migrations run up and down;
- query builder always parameterizes inputs;
- transactions roll back on handler failure;
- generated apps can reset test database state.

### Container And IaC

Purpose: keep local development close to production failure shapes.

Future checks:

- generated Compose file validates;
- web and API containers build from a clean checkout;
- web installs with Bun from `bun.lock`;
- Tailwind is a required generated web dependency and build step;
- web builds content-hashed React JS and CSS assets;
- web serves the HTML shell and asset manifest with `no-store`;
- web serves hashed assets with one-year immutable cache headers;
- web, API, and Postgres run as separate services;
- health checks converge;
- generated API logs the external web URL used for API proxying;
- login fails before the first user is registered;
- registration through `/api/register` creates the first user, sets a cookie,
  and returns JSON;
- login through `/api/login` works after registration;
- generated Compose config declares file-watch rebuilds for `web`
  source, web package/config files, `api/`, `db/`, and
  `api/Dockerfile` changes;
- generated apps include an env contract in `carbide.toml` and
  `.env.example`, but do not scaffold `README.md`, `AGENTS.md`, or `agents.d/`;
- generated frontend structure keeps L1/L2/L3 as class-ownership layers
  reflected in component directories;
- generated apps include a Bun/React/Tailwind web container, Go API
  container, and Postgres db container;
- Bun web service proxies `/api` and `/health` to the API service;
- `/api/me` reports anonymous and authenticated state correctly;
- `/dashboard` is served by the React app shell;
- `/dashboard` and `/` reference content-hashed JS and CSS assets;
- restart behavior preserves Postgres data;
- environment contract rejects missing required values;
- secret values are never printed by `carbide health env`;
- `carbide health` verifies generated project shape, required config,
  env/secrets, Compose, and legacy-regression markers without starting
  containers;
- `carbide health runtime` runs the Docker-backed health/auth/dashboard flow
  and stops containers it started;
- browser-exposed variables cannot be marked secret;
- framework-owned keys are visible in the contract and protected from casual app
  override;
- `carbide deploy check prod` classifies missing, invalid, preview-only, and
  apply-supported deploy targets;
- `carbide deploy preview prod` is non-mutating and reports either the
  checked-in production target plan, including environment hosts and roles, or
  the missing-target state;
- `carbide deploy apply prod` refuses unknown targets and runs only when `prod`
  is a checked-in deploy target with implemented apply semantics;
- `ssh-compose-environment` targets are previewable and validated while
  clustered apply remains guarded.

### CLI Golden Tests

Purpose: protect developer experience and generated files.

Future checks:

- `carbide new` creates the canonical directory structure;
- `carbide init` succeeds only in an empty directory;
- `carbide run dev` prints a compact startup summary and suppresses noisy
  Compose build output by default;
- `carbide run dev` prints only the working app/API URLs before the log stream,
  with no port-busy, demo-login, mode, status, stop, or watch-summary rows;
- `carbide run dev` shows full-width TTY-only per-container startup animation while
  Compose starts containers, without leaking progress control text into
  captured output, and without treating `NO_COLOR` as a request to disable
  terminal animation;
- `Ctrl+C` during `carbide run dev` detaches from live logs without running
  `docker compose down`;
- `carbide clean dev` normalizes session state without deleting volumes and
  runs `docker compose down --remove-orphans` when containers exist;
- `carbide stop dev` is the explicit teardown path, runs `docker compose down`,
  and shows full-width TTY-only per-container shutdown animation;
- CLI success, error, version, upgrade, and dev-stack output use the shared
  aligned renderer instead of scattered raw prints;
- `carbide run dev` streams web, API, db, and watch output
  through timestamped service-tagged rows after the stack is ready;
- `carbide status` prints a stable table of services, container names,
  published host ports, internal container ports, and status;
- `carbide urls json`, `carbide status json`, `carbide health json`,
  `carbide health env json`, `carbide health runtime json`,
  `carbide health framework json`, and deploy JSON subcommands emit valid,
  ANSI-free machine-readable state;
- `carbide audit` starts a Codex audit session, creates an audit workspace
  with a starter-reference snapshot and an audit brief, auto-launches Codex
  only in interactive terminals, and does not copy generated local artifacts
  or rewrite app code;
- `carbide health` prints a stable table of app-law checks;
- `carbide health env` validates the generated environment contract without
  printing secret values;
- `carbide health framework` runs source-repo regressions: shell syntax, Go
  CLI tests, repo contract, scaffold checks, and Docker smoke;
- `carbide deploy preview prod` prints the non-mutating deploy plan;
- `carbide deploy check prod` prints a non-mutating deploy classifier;
- `carbide deploy apply prod` applies only a checked-in production target and
  remains guarded for unknown targets and preview-only environment targets;
- `carbide follow logs` reattaches to live container logs and preserves
  timestamped, service-tagged rendering;
- `carbide run dev` writes `.carbide/log/dev.jsonl`, and `carbide logs` can
  query it by service, text, limit, and JSON output;
- generated apps contain no seeded demo account or demo credentials;
- generated files are deterministic;
- invalid commands print actionable errors;
- scaffolded apps pass the same CI checks as framework examples.

### Security Regression

Purpose: keep safe defaults from silently weakening.

Future checks:

- signed cookie tamper tests;
- CSRF token validation;
- secure cookie flags in production mode;
- upload size and type limits;
- SQL parameterization tests;
- secret values never printed in logs.

### Documentation Regression

Purpose: make the documentation site track framework behavior.

Initial checks:

- static docs site exists;
- custom domain is present;
- engineering plans are checked in.

Future checks:

- docs examples compile;
- generated command snippets are paste-ready;
- links resolve inside the docs site;
- release docs are versioned.
- framework runtime baselines are documented in
  `docs/engineering/VERSION_POLICY.md`, recorded in generated `carbide.toml`,
  and audited without automatic mutation.

## Current Implemented Checks

The first implemented CI job is intentionally small:

```sh
bash -n tests/contract/check_repo_contract.sh tests/scaffold/cli_scaffold.sh tests/smoke/starter_docker_flow.sh cli/bin/carbide cli/install.sh
(cd cli && go test ./...)
bash tests/contract/check_repo_contract.sh
bash tests/scaffold/cli_scaffold.sh
bash tests/smoke/starter_docker_flow.sh
carbide health framework
```

This protects the repo, generated starter, Docker dev topology, and
documentation deployment while the framework code is still being designed.
