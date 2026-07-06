# Carbide Docs App

This is the Carbide app that serves the checked-in documentation website from
`../site`.

Agent startup guidance is centralized at:

```text
https://carbide.ryangerardwilson.com/for/agents
```

Use this app only for the docs runtime and deploy loop:

```sh
carbide doctor
carbide deploy preview de-sci
carbide deploy apply de-sci
```
