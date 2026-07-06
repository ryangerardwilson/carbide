# Guiding Your Agents to Get Started

```shell
curl -fsSL https://raw.githubusercontent.com/ryangerardwilson/carbide/main/cli/install.sh | bash
carbide new demo
cd demo
carbide run dev
carbide doctor
carbide status
```

The installer builds the Go CLI. Human app names are accepted:
`carbide new "My Carbide App"` creates `my-carbide-app` while storing the
display name as `My Carbide App`.

`carbide help` prints the command reference. `carbide upgrade` updates the
installed CLI when a newer GitHub commit is available.
