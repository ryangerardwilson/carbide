# Scaffold Contract

`scaffold/` is the canonical generated Carbide app. Changes here affect every
future `carbide new` and `carbide init`.

## Generated App Shape

Generated app root:

```text
.
|-- api/
|-- db/
|-- web/
|-- carbide.toml
|-- docker-compose.yml
`-- .env.example
```

Every root directory is a standalone Docker service.

Carbide does not scaffold `README.md`, `AGENTS.md`, or `agents.d/`. If an app
owner creates local prose files later, those are app-owned context rather than
framework-owned contract files.

Generated `carbide.toml` records app identity, env contract, deploy targets,
and current starter runtime defaults.

After scaffold, the app owns its own code immediately. Carbide does not treat
existing app files as framework-managed.

## Web Contract

- TypeScript only under `web/src`.
- Bun owns install/build/server inside the web container.
- React owns browser state and auth/dashboard views.
- Tailwind is mandatory.
- `web/src/styles.css` must stay small: Tailwind import, `@source` globs, and
  `@custom-variant dark`.
- Do not put `html`, `body`, root font-size, width, line-height, layout, or
  component styling rules in `styles.css`; those belong in Tailwind utility
  classes and component tokens.
- Generated colors and light/dark variants belong in component Tailwind
  classes, especially `web/src/component/l1/tokens.ts`.
- Built-in scrollbar styling belongs in Tailwind tokens/classes using
  `scrollbar-width` and `scrollbar-color`, with light/dark variants. Do not add
  scrollbar pseudo-selector CSS to `styles.css`.
- Do not reintroduce generated `--carbide-*` color variables or `@theme` into
  `styles.css`.
- If a real app intentionally needs global product CSS, create
  `web/src/product.css`, import it explicitly, document the reason in local
  app docs if they exist, and update the law contract. Do not hide product CSS
  inside `web/src/styles.css`.

## API And DB Contract

- `api/` owns Go HTTP routes, auth, sessions, validation, JSON responses, and
  request logs.
- `db/` owns Postgres-facing helpers and migration state.
- Postgres is mandatory, not one adapter among many.

## Contract Changes

Any scaffold contract change should update:

- `tests/contract/check_repo_contract.sh`,
- `tests/scaffold/cli_scaffold.sh`,
- `cli/internal/cli/cli_test.go` when CLI behavior changes,
- relevant docs under `docs/engineering/`,
- docs site HTML when public docs copy changes.
