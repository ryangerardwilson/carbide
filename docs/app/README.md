# Carbide Docs App

This is the Carbide app that serves the checked-in documentation website from
`../site`.

Agent startup guidance is centralized at:

```text
https://carbide.ryangerardwilson.com/for/agents
```

Use this app only for the docs runtime and deploy loop:

```sh
export CARBIDE_DOCS_DEPLOY_SSH=<ssh-destination>
carbide doctor
carbide deploy check prod
carbide deploy preview prod
carbide deploy apply prod
```

Set `CARBIDE_DOCS_DEPLOY_SSH` in the shell or CI secret store before deploy.
