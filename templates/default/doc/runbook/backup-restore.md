# Backup And Restore

Postgres owns durable application state in a Carbide app.

The local development stack stores Postgres data in the `carbide_pgdata`
Compose volume. That volume is disposable for development, but production
targets must define explicit backup and restore behavior before they are
considered supported.

## Local Reset

To stop containers and remove the local database volume:

```sh
docker compose down -v --remove-orphans
```

## Production Rule

A production deploy target is incomplete until it documents:

- where backups are stored;
- how often backups run;
- how restore is tested;
- how migrations are ordered relative to deploys;
- who can access database credentials.
