# Taste Guide

This file describes Carbide's current starter taste.

Taste is allowed to evolve. It influences scaffold shape, starter UX, package
choices, runtime pins, UI organization, docs examples, and helper defaults.
Taste is not the same thing as law.

## Current Taste Areas

- React, Bun, Tailwind, Go, and Postgres as the current starter stack.
- Current pinned runtime images and package versions in scaffolded files.
- The generated auth, registration, dashboard, and theme behavior.
- Current frontend component organization and Tailwind usage style.
- Current docs wording and example flows.
- Current CLI presentation and audit/reporting format.

## Enforcement Rule

Generated apps may diverge from this taste.

`carbide health` should not fail only because an app no longer matches current
starter taste. Taste is for scaffolding and for audits driven by Codex with
user intent.
