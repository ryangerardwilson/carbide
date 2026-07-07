# Laws

These are the Carbide laws.

They are the part that should not drift across starter revisions, CLI
releases, or current framework taste. `carbide health` for generated apps
checks these laws and nothing broader.

## Generated App Laws

Agents should cite these clauses as `Law 1` through `Law 7`.

### Law 1. One App Repo

Every Carbide app is one repo.

### Law 2. Root Runtime Directories

The root runtime directories are `web/`, `api/`, and `db/`.

### Law 3. Checked-In Runtime Contracts

`carbide.toml` and `docker-compose.yml` are checked in.

### Law 4. Same-Origin Browser Flow

The browser entrypoint is `web`, and browser API traffic stays same-origin
through `/api` to `api`.

### Law 5. Postgres Is Required

Postgres is the required durable database.

### Law 6. Preview Before Apply

Deploy remains preview-before-apply.

### Law 7. Secrets Are Never Printed

Carbide output, docs, and agent guidance never print secrets.

## 2. Ownership Rule

`carbide new` and `carbide init` create the starter.

After scaffold, the app owns its own code immediately. Carbide does not treat
existing app files as framework-managed and does not rewrite them as part of
upgrade, migration, or health checks. Carbide does not scaffold or validate
`README.md`, `AGENTS.md`, or `agents.d/`; if an app owner creates local prose
files later, they are app-owned context.

## 3. Audit Rule

`carbide audit` may create local comparison material under `.carbide/`, but it
does not rewrite app code.

If app code changes during an audit, those changes are Codex edits applied
intentionally inside the app because the user wants them.
