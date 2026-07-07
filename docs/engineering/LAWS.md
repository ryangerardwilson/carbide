# Laws

These are the Carbide laws.

They are the part that should not drift across starter revisions, CLI
releases, or current framework taste. `carbide health` for generated apps
checks these laws and nothing broader.

## Generated App Laws

1. One app repo.
2. Root runtime directories are `web/`, `api/`, and `db/`.
3. `carbide.toml` and `docker-compose.yml` are checked in.
4. The browser entrypoint is `web`, and browser API traffic stays same-origin
   through `/api` to `api`.
5. Postgres is the required durable database.
6. Deploy remains preview-before-apply.
7. Carbide output, docs, and agent guidance never print secrets.

## Ownership Rule

`carbide new` and `carbide init` create the starter.

After scaffold, the app owns its own code immediately. Carbide does not treat
existing app files as framework-managed and does not rewrite them as part of
upgrade, migration, or health checks. Carbide does not scaffold or validate
`README.md`, `AGENTS.md`, or `agents.d/`; if an app owner creates local prose
files later, they are app-owned context.

## Audit Rule

`carbide audit` may create local comparison material under `.carbide/`, but it
does not rewrite app code.

If app code changes during an audit, those changes are Codex edits applied
intentionally inside the app because the user wants them.
