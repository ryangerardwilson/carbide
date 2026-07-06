# Regression Checks

Use the narrowest useful check while editing, then broader checks before
shipping.

## Fast Checks

```sh
cd cli && go test ./...
bash tests/contract/check_repo_contract.sh
PATH=/home/ryan/.local/share/mise/installs/go/1.26.4/bin:$PATH bash tests/scaffold/cli_scaffold.sh
```

## Framework Check

```sh
carbide health framework
```

This runs shell syntax, Go CLI tests, repo contract, scaffold checks, and the
generated Docker smoke flow.

## Web Checks

For scaffold web:

```sh
cd scaffold/web
bun install --frozen-lockfile
bun run typecheck
bun run assets:build
```

For docs web:

```sh
cd docs/app/web
bun run typecheck
bun run assets:build
cd ../
docker compose build web
```

## Health Checks

```sh
cd scaffold && carbide health
cd docs/app && carbide health
```

Run `carbide health runtime` when container behavior changed and Docker is
available.

## Docs Agent Route

After deploying docs, verify the public agent guide:

```sh
bash tests/smoke/docs_for_agents_http.sh
```
