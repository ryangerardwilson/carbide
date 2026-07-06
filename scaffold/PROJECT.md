# __PROJECT_NAME__

This file is for app-specific product truth. Keep framework setup and agent
workflow instructions in `AGENTS.md` and the central Carbide agent guide.

## Product Truth

- This app is named `__PROJECT_NAME__`.
- Replace this section with the product's domain, promise, and core workflow.

## Users And Roles

- First user: created through the browser registration flow.
- Add real user roles here when the product needs them.

## Business Rules

- Keep durable product invariants here.
- Keep secrets, runtime configuration, and deploy targets in `carbide.toml` and
  the deployment layer, not in this file.

## Acceptance Criteria

- The app starts with `carbide run dev`.
- The project contract passes with `carbide doctor`.
- Runtime behavior that touches containers or auth is verified with
  `carbide doctor runtime` when Docker is available.
