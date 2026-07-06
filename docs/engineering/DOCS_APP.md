# Docs App

The public documentation site is itself a Carbide app in `docs/app`.

Public URL:

```text
https://carbide.ryangerardwilson.com
```

## Ownership

- `docs/site/`: checked-in static documentation pages.
- `docs/app/web/`: Bun/TypeScript server that serves `docs/site`, rewrites
  docs HTML through React component contracts, and serves generated browser
  assets.
- `docs/app/api/`: small Go health/version API.
- `docs/app/db/`: docs app Postgres migration state.

## Rules

- Docs app should dogfood the current scaffold frontend contract unless there
  is a clear reason not to.
- `docs/app/web/src/styles.css` must match `scaffold/web/src/styles.css`.
- Edit durable docs content in `docs/engineering/` and checked-in public pages
  under `docs/site/` as appropriate.
- Run `bun run assets:build` in `docs/app/web` after docs style changes so
  `docs/site/assets/styles.css` is regenerated.
- Set `CARBIDE_DOCS_DEPLOY_SSH` from shell env or CI secrets.
- Deploy with `carbide deploy apply prod` from `docs/app` only after local
  checks pass.
