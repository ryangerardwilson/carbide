# Deployment

Carbide deploys from checked-in targets in `carbide.toml`. A target is an
environment, not just a machine. The simplest environment is one VM running
Docker Compose. Larger environments declare hosts and roles so the topology is
reviewable before anything mutates.

Check first, then preview, then apply:

```sh
carbide deploy check prod
carbide deploy preview prod
carbide deploy apply prod
```

`prod` is the canonical production target name. `carbide deploy check prod`
classifies the target as `missing-target`, `preview-only`, `invalid-config`, or
`apply-supported`. `carbide deploy apply prod` mutates infrastructure only when
`prod` is a checked-in target with implemented apply semantics.

Carbide supports `ssh-compose` apply for a checked-in single-VM target. New
generated apps ship with no deploy target, so `apply prod` refuses until one
exists. `ssh-compose-environment` validates and previews multi-VM topology, but
apply is guarded until clustered orchestration is implemented.

Agents and CI should use:

```sh
carbide deploy check prod json
carbide deploy preview prod json
```

## Single VM

Use `type = "ssh-compose"` when one VM runs the `web`, `api`, and `db`
containers through Docker Compose.

```toml
[deploy]
preview_before_apply = true

[deploy.hosts.prod]
ssh = "${CARBIDE_DEPLOY_SSH}"

[deploy.targets.prod]
type = "ssh-compose"
host = "prod"
domain = "app.example.com"
remote_path = "/opt/my-carbide-app"
source_path = "."
compose_file = "docker-compose.yml"
project_directory = "."
public_port = 18080
health_path = "/health"
nginx = true
nginx_site = "my-carbide-app"
```

`host` can be a direct SSH destination or the name of a checked-in
`[deploy.hosts.*]` entry. Use the host table when the SSH destination should
come from shell env or CI secrets instead of a repo-local alias.

The remote VM needs Docker Compose available. If `nginx = true`, the remote user
must be able to run `sudo`; Carbide writes an nginx site for `domain`, proxies
to `public_port`, and uses Certbot when a certificate is not already present.

Deploy with:

```sh
carbide deploy check prod
carbide deploy preview prod
carbide deploy apply prod
```

On apply, Carbide syncs `source_path` to `remote_path`, creates `.env` on the
remote host if it does not exist, runs Compose config, starts the stack, updates
nginx when requested, and checks `health_path` through `127.0.0.1:public_port`.

## Multiple VMs

Use `type = "ssh-compose-environment"` when the environment has explicit hosts
and roles. This is the Carbide shape for multi-server apps.

```toml
[deploy]
preview_before_apply = true

[deploy.hosts.web-1]
ssh = "web-1"
description = "Public web entrypoint."

[deploy.hosts.api-1]
ssh = "api-1"
description = "Private API host."

[deploy.hosts.db-1]
ssh = "db-1"
description = "Primary Postgres host."

[deploy.targets.prod]
type = "ssh-compose-environment"
domain = "app.example.com"
remote_path = "/opt/my-carbide-app"
source_path = "."
compose_file = "docker-compose.yml"
project_directory = "."
health_path = "/health"
strategy = "preview-only"

[deploy.targets.prod.roles.web]
hosts = ["web-1"]
public_port = 18080
nginx = true

[deploy.targets.prod.roles.api]
hosts = ["api-1"]

[deploy.targets.prod.roles.db]
hosts = ["db-1"]
primary = "db-1"
migration = "once"
```

Preview with:

```sh
carbide deploy check prod
carbide deploy preview prod
```

The preview validates that each role references known hosts, that `web`, `api`,
and `db` roles exist, that the `db` role has a primary host, and that migrations
are declared as `once`.

`carbide deploy apply prod` is intentionally guarded for this target type
until Carbide implements clustered orchestration, migration ordering, health
gates, load-balancer changes, and rollback behavior.

## Target Fields

- `type`: `ssh-compose` for one VM, `ssh-compose-environment` for multi-VM
  topology.
- `host`: direct SSH destination for a single-VM `ssh-compose` target, or the
  name of a matching `[deploy.hosts.*]` entry.
- `domain`: public DNS name for the environment.
- `remote_path`: absolute path on the remote host.
- `source_path`: local path to sync before deploy.
- `compose_file`: Compose file path relative to `remote_path`.
- `project_directory`: Compose project directory relative to `remote_path`.
- `public_port`: host port exposed by the web service.
- `health_path`: HTTP path checked after deploy.
- `nginx`: whether Carbide should manage nginx for the public entrypoint.
- `nginx_site`: lowercase nginx site name.
- `deploy.hosts.*.ssh`: SSH destination for a named multi-VM host.
- `roles.*.hosts`: host names assigned to a role.
- `roles.db.primary`: primary Postgres host.
- `roles.db.migration`: currently must be `once`.
