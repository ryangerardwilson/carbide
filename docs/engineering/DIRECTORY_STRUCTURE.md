# Directory Structure

Carbide uses a narrow root so CLI code, generated app scaffolding, tests, and
documentation have clear ownership.

```text
.
|-- .github/
|   `-- workflows/
|       |-- ci.yml
|       `-- pages.yml
|-- cli/
|   |-- bin/
|   |   `-- carbide
|   |-- cmd/
|   |   `-- carbide/
|   |       `-- main.go
|   |-- go.mod
|   |-- install.sh
|   `-- internal/
|       `-- cli/
|           |-- cli.go
|           `-- cli_test.go
|-- docs/
|   |-- engineering/
|   |   |-- CI_CD_REGRESSION_TESTS.md
|   |   |-- CREATE_YOUR_FIRST_APP.md
|   |   |-- DEPLOYMENT.md
|   |   |-- DIRECTORY_STRUCTURE.md
|   |   |-- FRONTEND_STARTER_CONTRACT.md
|   |   `-- VERSION_POLICY.md
|   `-- site/
|       |-- CNAME
|       |-- assets/
|       |   `-- styles.css
|       |-- ci-cd-regression-tests.html
|       |-- create-your-first-app.html
|       |-- deployment.html
|       |-- frontend-starter-contract.html
|       |-- index.html
|       |-- repo-structure.html
|       `-- version-policy.html
|-- scaffold/
|   |-- AGENTS.md
|   |-- README.md
|   |-- agents.d/
|   |   |-- BACKUP_RESTORE.md
|   |   |-- DEPLOY.md
|   |   |-- ENVIRONMENT.md
|   |   `-- TAILWIND_COMPONENTS.md
|   |-- api/
|   |   |-- auth.go
|   |   |-- Dockerfile
|   |   |-- go.mod
|   |   |-- go.sum
|   |   |-- main.go
|   |   `-- routes.go
|   |-- carbide.toml
|   |-- db/
|   |   |-- go.mod
|   |   |-- go.sum
|   |   |-- migration/
|   |   |   `-- 001_auth.sql
|   |   |-- session.go
|   |   `-- user.go
|   |-- docker-compose.yml
|   `-- web/
|       |-- Dockerfile
|       |-- bun.lock
|       |-- index.html
|       |-- package.json
|       `-- src/
|           |-- component/
|           |-- lib/
|           |-- main.jsx
|           |-- server.jsx
|           |-- write-index.mjs
|           `-- styles.css
|-- tests/
|   |-- contract/
|   |   `-- check_repo_contract.sh
|   |-- scaffold/
|   |   `-- cli_scaffold.sh
|   `-- smoke/
|       `-- starter_docker_flow.sh
`-- README.md
```

## Ownership

- `.github/workflows/`: CI and documentation deployment.
- `cli/bin/carbide`: source checkout launcher for the Go CLI.
- `cli/cmd/carbide/`: installable CLI entrypoint.
- `cli/go.mod`: Go module definition for the CLI.
- `cli/internal/cli/`: Go implementation of the CLI and its unit tests.
- `docs/engineering/`: source-of-truth engineering plans.
- `docs/site/`: static documentation site served by the Carbide docs app.
- `scaffold/`: generated app source exactly as `carbide new` and
  `carbide init` write it to disk.
- At the generated scaffold root, every directory except `agents.d/` maps to a
  standalone Docker service: `web/`, `api/`, and `db/`.
- `scaffold/AGENTS.md`: generated agent-facing entrypoint for runtime shape,
  operating context, and safe default commands.
- `scaffold/api/`: generated Go HTTP/API server, auth, routing, session, and
  JSON response code, including its Go module and API Dockerfile.
- `scaffold/carbide.toml`: generated project metadata, default dev port,
  required/optional/secret/browser-exposed/framework-owned environment
  contract, and deploy guardrails.
- `scaffold/db/`: generated Postgres-backed data access code and its Go
  module.
- `scaffold/db/migration/`: checked-in generated schema state.
- `scaffold/agents.d/`: generated local operating notes for
  environment, deploy preview/apply, backup/restore, and Tailwind component
  organization.
- `scaffold/web/`: generated Bun/React/Tailwind web app,
  web container source, browser UI, and same-origin API proxy.
- `scaffold/web/src/component/l1/`: generated primitive UI elements and
  Tailwind utility tokens.
- `scaffold/web/src/component/l2/`: generated composed UI patterns such as
  forms and app layouts.
- `scaffold/web/src/component/l3/`: generated auth, dashboard, and loading
  screens.
- `scaffold/web/src/lib/`: generated non-component browser helpers such as the
  `cx()` class-name helper.
- `scaffold/docker-compose.yml`: generated local Docker Compose infrastructure.
- `tests/contract/`: repository contract checks.
- `tests/scaffold/`: generated-project and CLI scaffold checks.
- `tests/smoke/`: end-to-end smoke checks that boot generated apps.
- `cli/install.sh`: GitHub URL installer that builds the Go CLI and places
  `carbide` on the user's PATH.

## First Implementation Rule

Empty directories are not kept as placeholders. When a directory gains
behavior, its first file should make that behavior testable. The framework
root intentionally has no placeholder `src/`, `infra/`, or `examples/`
directory. In generated apps, avoid non-container root directories beside
`agents.d`; shared runtime coordination belongs in `docker-compose.yml`.
