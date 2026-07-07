# Frontend Starter Contract

Carbide's default app uses a Bun/React/Tailwind web service, Go API backend,
and Postgres database. The web service is a mandatory Bun container in the
default local topology, not a host-installed JavaScript tooling requirement.

## Product Decision

The default Carbide UI should be React served by Bun, not a custom template
language and not an opinionated design-system framework.

This keeps frontend authoring inside a mature ecosystem while preserving the
core Carbide bet: elegant monorepo structure, Docker-first local development,
Go API logic, Postgres state, and explicit environment/deploy
contracts.

## Runtime Model

```text
browser -> web container -> /api proxy -> api Go container -> db
```

- `web` owns Bun, React, Tailwind, browser routes, forms, dashboard UI,
  and the same-origin proxy.
- `api` owns Go API routes, auth, session cookies, validation, and JSON.
- `db` owns durable Postgres state.

The web service is the public entrypoint. It proxies `/api` and `/health` to the
API service so browser requests stay same-origin.

## Authoring Model

Generated apps start with:

```text
web/
|-- Dockerfile
|-- bun.lock
|-- index.html
|-- package.json
|-- tsconfig.json
`-- src/
    |-- component/
    |   |-- l1/
    |   |   |-- Button.tsx
    |   |   |-- Field.tsx
    |   |   |-- Surface.tsx
    |   |   |-- Text.tsx
    |   |   |-- ThemeToggle.tsx
    |   |   |-- index.ts
    |   |   `-- tokens.ts
    |   |-- l2/
    |   |   |-- AuthForm.tsx
    |   |   |-- Layouts.tsx
    |   |   `-- index.ts
    |   |-- l3/
    |   |   |-- AuthView.tsx
    |   |   |-- DashboardView.tsx
    |   |   |-- LoadingView.tsx
    |   |   `-- index.ts
    |-- lib/
    |   |-- cx.ts
    |   `-- types.ts
    |-- main.tsx
    |-- server.ts
    |-- styles.d.ts
    |-- write-index.ts
    `-- styles.css
```

Generated apps place the web app under `web/` so the project mirrors
runtime boundaries. HTTP code and its Go module live under `api/`; Postgres
access, its Go module, and migrations live under `db/`.

The default UI uses TypeScript and includes register, login, logout,
dashboard, and a conventional left-sidebar app shell. React components call
same-origin `/api` endpoints with `credentials: "include"` so the API can own
HttpOnly cookies.

## Component Stance

Carbide scaffolds a visible L1/L2/L3 component hierarchy so the Tailwind class
organization is mirrored in the project tree:

- `component/l1`: primitives and Tailwind utility tokens.
- `component/l2`: reusable composed UI patterns and layouts.
- `component/l3`: product screens and domain-specific sections.

That hierarchy is part of the generated starter contract. App teams may evolve
it later, but `carbide new` should teach the intended organization from the
first generated project.

Carbide's component-design taste should track Tailwind Plus / Catalyst more
than bespoke framework chrome:

- default to Application UI patterns for product work;
- keep components production-ready, fully responsive, accessible, and easy to
  customize;
- keep utility classes in component markup instead of pushing design into a
  separate CSS abstraction layer;
- favor neutral palettes, clear hierarchy, dense but readable spacing, and
  operational layouts over decorative marketing composition;
- include normal interaction and data states in the component contract:
  default, hover, focus-visible, active, disabled, loading, empty, and error
  where relevant;
- keep component APIs small and composable, usually around content, intent,
  size, state, and `className`;
- rely on conventional application primitives such as sidebars, stacked
  headers, form sections, cards, tables, dialogs, dropdowns, and stats rows
  before inventing a new surface pattern.

## Styling

Generated apps use Tailwind as the mandatory styling path. `styles.css` is the
Tailwind input file, and the container builds generated CSS with the checked-in
Bun lockfile.

`styles.css` contains only the Tailwind import, TypeScript-aware source
directives, and the `data-theme` dark variant. It does not own global
`html`/`body` sizing, layout defaults, or a generated color-variable palette.
`tokens.ts` contains reusable Tailwind utility groups for the generated auth
and dashboard UI, including starter light/dark visual choices and built-in
scrollbar utilities. The scaffold does not add a parallel `theme.css` file.
`carbide health` rejects global `html`/`body` sizing, custom class selectors,
ID selectors, scrollbar pseudo-selectors, `--carbide-*` color variables,
`@theme`, `@apply`, `@layer`, keyframes, media rules, and container rules in
`styles.css`; those belong in Tailwind utility classes and component class
layers.

Theme mode follows one starter pattern:

- `ThemeToggle.tsx` lives in `component/l1/`.
- Theme state stays browser-local by default.
- The selected mode is stored in `localStorage`.
- `system` resolves through `matchMedia('(prefers-color-scheme: dark)')`.
- The app sets `document.documentElement.dataset.theme` to the resolved
  `light` or `dark` value.
- The app sets `document.documentElement.dataset.themeMode` to the selected
  mode and mirrors the resolved mode with `document.documentElement.style.colorScheme`.
- `styles.css` keeps Tailwind's dark variant bound to
  `[data-theme="dark"]`.
- L2 and L3 components receive `mode`, `resolved`, and `onMode`; they should
  not each create their own theme state.

Use these sanctioned paths instead:

- put reusable class groups in `web/src/component/l1/tokens.ts`;
- compose component variants in TypeScript with `cx()`;
- keep responsive layout, typography, spacing, color, and scrollbars in
  Tailwind utilities;
- keep third-party CSS imports explicit and product-owned;
- if a real product intentionally needs global CSS, create
  `web/src/product.css`, import it explicitly, document the reason in local
  app docs if they exist, and update the law contract.

`typecheck` runs `tsc --noEmit`. Docker builds run typecheck before building
browser assets, so broken component props, API response shapes, or Bun server
types fail before the container starts.

## Browser Asset Contract

Generated apps cache-bust React browser assets by default. `assets:build`
runs Tailwind, runs Bun's browser build with content-hashed filenames under
`web/public/assets/`, and writes `web/public/index.html` plus
`web/public/asset-manifest.json`.

The HTML app shell and asset manifest are served with `Cache-Control: no-store`
so browsers always learn the current hashed asset names. Files under
`/assets/` are served with `Cache-Control: public, max-age=31536000, immutable`
because their URLs change when content changes.

`web/public/` and `web/src/tailwind.css` are generated build output and are
ignored in generated app repositories.

## Regression Tests

The frontend contract needs dedicated regression coverage:

- generated apps include a Bun/React/Tailwind web container;
- generated apps include a Go API container;
- generated apps include a Postgres db container;
- Bun web service proxies `/api` and `/health` to the API service;
- auth uses same-origin cookies without CORS setup;
- login returns JSON and sets a session cookie;
- `/api/me` reports authenticated and anonymous states correctly;
- `/dashboard` is served by the React app shell;
- the React app shell references content-hashed JS and CSS assets;
- the app shell and asset manifest are served `no-store`;
- hashed assets are served with one-year immutable caching;
- web, API, and db watch paths are present in Compose;
- generated web service installs with `bun install --frozen-lockfile` and builds
  with `bun run build`;
- generated web service runs `tsc --noEmit` through `bun run typecheck`;
- Tailwind is present and required in the generated web service;
- generated web app includes `component/l1`, `component/l2`, and `component/l3`
  starter tiers;
- generated `main.tsx` imports starter screens and does not reimplement
  dashboard or auth screen markup inline;
- generated apps keep primitive UI in L1, reusable patterns in L2, and product
  screens in L3.
