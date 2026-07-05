# Deploy Preview And Apply

Carbide uses a preview-before-apply rule for infrastructure.

Preview means: show what would change.

Apply means: actually change infrastructure.

Until this app has a real deploy target, `carbide deploy preview <target>`
prints a non-mutating plan and `carbide deploy apply <target>` refuses to run.

Deploy targets should be modeled as environments. A small app can use one
`ssh-compose` host. A larger app should define `deploy.hosts.*` and an
`ssh-compose-environment` target with `web`, `api`, and `db` roles. Preview is
allowed for both shapes; clustered apply remains guarded until the app has
explicit migration, health, load-balancer, and rollback rules.

## Intended Flow

```sh
carbide deploy preview dev
carbide deploy apply dev
```

Production must be stricter:

```sh
carbide deploy preview prod
carbide deploy apply prod
```

The production apply path must require explicit confirmation once a production
target exists.

## Rules

- Never mutate infrastructure without a preview.
- Prefer provider identity such as GitHub OIDC over long-lived cloud keys.
- Keep deploy targets checked in, reviewable, and recoverable.
- Store real secrets in the deploy/IaC layer, not in source code.
- Treat a target as an environment made of hosts and roles, not merely a single
  host string.
