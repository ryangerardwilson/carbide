# Carbide Agent Guide

This file is for agents and maintainers iterating on the Carbide framework
repo. The public README is for humans evaluating and installing Carbide; do not
turn it into a roadmap, scratchpad, or internal operating manual.

## Source Of Truth

- Product promise and README boundary: `agents.d/PRODUCT_CONTRACT.md`
- Repo layout and ownership: `agents.d/REPO_STRUCTURE.md`
- CLI and versioning rules: `agents.d/CLI_AND_VERSIONING.md`
- Generated app contract: `agents.d/SCAFFOLD_CONTRACT.md`
- Documentation app contract: `agents.d/DOCS_APP.md`
- Regression checks: `agents.d/REGRESSION_CHECKS.md`
- Roadmap and future work: `agents.d/ROADMAP.md`

Read the smallest relevant file before editing. Do not load every file by
default.

## Hard Rules

- Keep Carbide monorepo-first. Containers define runtime boundaries; they do
  not justify default microservice sprawl.
- Keep the generated app Docker-first: `web`, `api`, and `db` are separate
  containers in one app repo.
- Keep README human-first: what Carbide is, what the user gets, how to install,
  and where to read more.
- Put agent iteration rules, roadmap phases, internal contracts, and regression
  details in `AGENTS.md` or `agents.d/`.
- Keep root directories intentional. The framework repo root should stay small:
  `agents.d`, `cli`, `docs`, `scaffold`, and `tests`.
- Do not add root-level generated app files. Generated-app examples belong in
  `scaffold/`; docs-app runtime belongs in `docs/app/`.
- Before changing scaffold output, update contract tests and run the scaffold
  checks.

## Normal Verification

Use the narrowest check that proves the change, then run broader checks before
shipping contract changes:

```sh
cd cli && go test ./...
bash tests/contract/check_repo_contract.sh
PATH=/home/ryan/.local/share/mise/installs/go/1.26.4/bin:$PATH bash tests/scaffold/cli_scaffold.sh
carbide doctor framework
```

For docs app changes:

```sh
cd docs/app/web && bun run typecheck && bun run assets:build
cd docs/app && docker compose build web
```

For scaffold web changes:

```sh
cd scaffold/web && bun install --frozen-lockfile && bun run typecheck && bun run assets:build
```

