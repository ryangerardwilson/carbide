# Frontend Contract

Carbide's default app uses a Bun/React/Tailwind frontend, Go API backend, and
Postgres database. The frontend is a mandatory Bun container in the default
local topology, not a host-installed JavaScript tooling requirement.

## Product Decision

The default Carbide UI should be React served by Bun, not a custom Blade-like
template system.

This keeps frontend authoring inside a mature ecosystem while preserving the
core Carbide bet: Go owns backend logic, auth, sessions, database access, and
the framework runtime contract.

## Runtime Model

```text
browser -> frontend container -> /api proxy -> backend Go container -> Postgres
```

- `frontend` owns Bun, React, Tailwind, browser routes, forms, dashboard UI,
  and the same-origin proxy.
- `backend` owns Go API routes, auth, session cookies, validation, and JSON.
- `db` owns durable Postgres state.

The frontend is the public entrypoint. It proxies `/api` and `/health` to the
backend so browser requests stay same-origin.

## Authoring Model

Generated apps start with:

```text
view/
`-- web/
    |-- Dockerfile
    |-- bun.lock
    |-- index.html
    |-- package.json
    `-- src/
        |-- component/
        |   |-- l1/
        |   |   |-- theme.css
        |   |   `-- tokens.js
        |   |-- l2/
        |   |-- l3/
        |   `-- utils.js
        |-- main.jsx
        |-- server.jsx
        `-- styles.css
```

Generated apps place the web app under `view/web/` so the project keeps the
MVC directory shape: `model/`, `view/`, and `controller/`.

The default UI includes register, login, logout, dashboard, and a built-in
component library. React components call same-origin `/api` endpoints with
`credentials: "include"` so the backend can own HttpOnly cookies.

## Component Contract

Generated apps enforce a three-level component structure:

- `component/l1/`: theme tokens, font and color variables, semantic UI
  classes, and design primitives such as buttons, fields, surfaces, text,
  badges, metrics, and code treatments. L1 never knows auth, routing, API
  calls, or app-specific state.
- `component/l2/`: reusable UX patterns composed from L1. This includes
  Lessons, Dropdown/Menu, Modal/Dialog/Slideover, Accordion/Disclosure,
  Carousel, Tabs, Notifications, Radio/Radio Group, Toggle/Switch, Tooltip,
  Popover, Listbox/Select, Combobox/Autocomplete, text editor adapters
  (Trix, Quill, SimpleMDE), chart adapters (Chart.js, ApexCharts), enhanced
  select adapters (Select2, Choices.js), calendar/date adapters (Flatpickr,
  Date Range Picker, FullCalendar), carousel adapters (Glide, Splide), and
  layout patterns for dashboards and landing pages. `DashboardLayout` is a
  conventional app shell with a left sidebar, section nav, account footer, and
  main work area.
- `component/l3/`: product and app surfaces composed from L2 patterns. The
  starter ships `AuthView`, `DashboardView`, `ComponentLibraryView`, and
  `LoadingView`.

`main.jsx` owns browser route state and API calls. It imports L3 views instead
of building screen markup inline. L3 views may pass product data downward; L2
and L1 stay reusable.

## Styling

Generated apps use Tailwind as the mandatory styling path. `styles.css` is the
Tailwind input file, and the container builds generated CSS with the checked-in
Bun lockfile.

The component library uses Tailwind layout classes directly, but font and color
scheme decisions live in L1. `theme.css` owns CSS custom properties and
semantic `cb-*` classes. `tokens.js` exports the same contract for React
components. Third-party integration names are represented as adapter
components that render useful built-in fallbacks without adding mandatory
frontend dependencies to every new app.

## Regression Tests

The frontend contract needs dedicated regression coverage:

- generated apps include a Bun/React/Tailwind frontend container;
- generated apps include a Go backend/API container;
- generated apps include a Postgres database container;
- Bun frontend proxies `/api` and `/health` to the backend;
- auth uses same-origin cookies without CORS setup;
- login returns JSON and sets a session cookie;
- `/api/me` reports authenticated and anonymous states correctly;
- `/dashboard` is served by the React app shell;
- frontend and backend watch paths are present in Compose;
- generated frontend installs with `bun install --frozen-lockfile` and builds
  with `bun run build`;
- Tailwind is present and required in the generated frontend;
- generated frontend includes `component/l1`, `component/l2`, and
  `component/l3`;
- generated L1 includes `theme.css` and `tokens.js`, and `styles.css` imports
  the L1 theme after Tailwind;
- generated `main.jsx` imports L3 views and does not reimplement dashboard or
  auth screen markup inline;
- L2 includes the Alpine-style interaction patterns and named integration
  adapter components.
