# Deploy Preview And Apply

Carbide uses a preview-before-apply rule for infrastructure.

Preview means: show what would change.

Apply means: actually change infrastructure.

Until this app has a real deploy target, `carbide deploy preview <target>`
prints a non-mutating plan and `carbide deploy apply <target>` refuses to run.

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
