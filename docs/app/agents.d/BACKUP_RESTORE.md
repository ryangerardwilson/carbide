# Backup And Restore

The docs content is checked into git under `docs/site`.

Postgres exists for runtime wiring and future docs metadata. If it begins to
own durable application state, back up the `carbide_docs_pgdata` volume before
host changes or destructive deploy work.
