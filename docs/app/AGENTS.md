# Carbide Docs App Agent Context

The central agent startup guide is:

```text
https://carbide.ryangerardwilson.com/for/agents
```

The checked-in source for that Markdown route is `../site/for/agents.md`.
This docs app intentionally does not include `agents.d/`.

## App Truth

This directory is a Carbide application used to deploy the checked-in
documentation website from `../site`.

- `web/` is the public entrypoint. It serves `../site`, rewrites docs HTML
  through Tailwind/React component contracts, and proxies `/api` and `/health`.
- `api/` exposes deploy health checks.
- `db/` owns docs app Postgres migration state.
- `carbide.toml` owns deploy targets and runtime/env contracts.

## Safe Commands

```sh
export CARBIDE_DOCS_DEPLOY_SSH=<ssh-destination>
carbide doctor
carbide deploy check prod
carbide deploy preview prod
carbide deploy apply prod
```

Set `CARBIDE_DOCS_DEPLOY_SSH` in the shell or CI secret store before deploy.
