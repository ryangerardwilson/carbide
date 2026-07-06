# CLI And Versioning

The `carbide` CLI is a compiled Go binary. The development wrapper at
`cli/bin/carbide` builds from source when `CARBIDE_HOME` points at this repo.
The installed binary normally lives at `.cli/bin/carbide` and is symlinked from
`~/.local/bin/carbide`.

## Version Rule

The CLI version string in `cli/internal/cli/cli.go` must be a real product
version unless the repo is intentionally between releases.

Current product version:

```text
0.1.0
```

When changing the version:

- update `cli/internal/cli/cli.go`,
- update generated `scaffold/carbide.toml`,
- update tests that assert command output,
- rebuild the installed binary if Ryan will use the command immediately.

## CLI Grammar

Use subject/action/object style where possible:

- `carbide clean dev`
- `carbide run dev`
- `carbide stop dev`
- `carbide follow logs`
- `carbide project migrate`

Do not add user-facing dash-flag aliases for core actions. `help`, `version`,
and `upgrade` are commands, not options.

Machine-readable output uses command-shaped JSON subcommands, not dash flags:

```sh
carbide urls json
carbide status json
carbide doctor json
carbide doctor env json
carbide doctor runtime json
carbide doctor framework json
carbide deploy check prod json
carbide deploy preview prod json
```

`-h` and `--help` may remain accepted as hidden compatibility aliases while
they exist, but do not show them in help text, README examples, or new command
docs.

## Installed Binary Refresh

After CLI source changes:

```sh
cd cli
commit="$(git rev-parse --short HEAD)"
tmp_bin="$(mktemp)"
go build -ldflags "-X github.com/ryangerardwilson/carbide/cli/internal/cli.commit=$commit" -o "$tmp_bin" ./cmd/carbide
chmod +x "$tmp_bin"
mv "$tmp_bin" /home/ryan/Apps/carbide/.cli/bin/carbide
ln -sfn /home/ryan/Apps/carbide/.cli/bin/carbide /home/ryan/.local/bin/carbide
```

Use `-dirty` in the commit metadata only when intentionally installing an
uncommitted working tree.
