# Repo Structure

The Carbide framework repo root is intentionally small:

```text
.
|-- cli/
|-- docs/
|-- scaffold/
|-- tests/
`-- README.md
```

## Ownership

- `cli/`: Go CLI source, installed CLI wrapper, parser, renderer, health,
  deploy, scaffold copy, logs, and tests.
- `docs/engineering/`: durable engineering docs.
- `docs/site/`: checked-in static docs pages served by the docs app.
- `docs/app/`: Carbide app that serves the public docs site.
- `scaffold/`: canonical generated app template.
- `tests/contract/`: repository shape and contract checks.
- `tests/scaffold/`: CLI-generated app checks.
- `tests/smoke/`: Docker-backed generated app smoke flow.
- `README.md`: framework-agent entrypoint and root operational guide.

## Root Rules

- No root `go.mod`; CLI Go module lives in `cli/`.
- No root `src`, `examples`, `infra`, `include`, or `templates`.
- No root `AGENTS.md`; framework-agent guidance lives in `README.md`.
- No root `agents.d`; `/for/agents` is the central agent startup surface.
- No generated-app source at root. Generated app shape lives in `scaffold/`.
- Do not duplicate scripts and tests with overlapping mandates. Contract
  checks live in `tests/contract`; scaffold checks live in `tests/scaffold`.
