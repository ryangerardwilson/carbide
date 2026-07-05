# Carbide Docs App

This directory is a Carbide application used to deploy the checked-in
documentation website from `../site`.

Use:

```sh
carbide doctor
carbide deploy preview de-sci
carbide deploy apply de-sci
```

The web container is the public entrypoint. It serves `docs/site` and proxies
`/api` and `/health` to the API container. Its source lives in `web/src`,
uses Tailwind, and keeps UI implementation boundaries in `component/l1`,
`component/l2`, and `component/l3`. The API container proves Postgres wiring
and exposes deploy health checks.
