# Deploy

Deploys use preview-before-apply.

`de-sci` is an `ssh-compose` target:

- Sync `docs/` to `/opt/carbide/docs`.
- Build and run `app/docker-compose.yml` on the remote host.
- Publish the web container on `127.0.0.1:18081`.
- Install an nginx HTTPS proxy for `carbide.ryangerardwilson.com`.
- Keep the nginx site name as `carbide`.

`de-sci-environment` is the same deploy expressed as an environment made of
hosts and roles. It is preview-only until Carbide implements clustered apply
semantics:

- `deploy.hosts.*` names SSH destinations.
- `roles.web` owns public entrypoints and nginx.
- `roles.api` owns API containers.
- `roles.db` owns the primary Postgres host and the migration-once rule.

The public DNS record points `carbide.ryangerardwilson.com` at the server, and
`/health` must pass through nginx before a deploy is considered complete.
