# Taste Guide

This file describes Carbide's current starter taste.

Taste is allowed to evolve. It influences scaffold shape, starter UX, package
choices, runtime pins, UI organization, docs examples, and helper defaults.
Taste is not the same thing as law.

## Current Taste Areas

- React, Bun, Tailwind, Go, and Postgres as the current starter stack.
- Current pinned runtime images and package versions in scaffolded files.
- The generated auth, registration, dashboard, and theme behavior.
- Current frontend component organization and Tailwind usage style, including
  `web/src/component/l1`, `l2`, and `l3` plus L1/L2/L3 class ownership inside
  reusable components.
- Current docs wording and example flows.
- Current CLI presentation and audit/reporting format.

## Enforcement Rule

Generated apps may diverge from this taste.

`carbide health` should not fail only because an app no longer matches current
starter taste. Taste is for scaffolding and for audits driven by Codex with
user intent.

## Current Frontend Taste

The current starter frontend teaches two related conventions:

- directory tiers:
  - `component/l1`: primitives and Tailwind utility tokens;
  - `component/l2`: reusable composed patterns and layouts;
  - `component/l3`: product screens and domain-specific sections.
- class ownership inside reusable components:
  - `l1`: structure and layout;
  - `l2`: geometry, spacing, borders, radii, and type scale;
  - `l3`: theme, color, state, motion, and interaction.

When a component is reusable or variant-heavy, prefer explicit class-layer
maps and `cx()` composition over long unreadable inline class strings.

Component design should also stay close to Tailwind Plus / Catalyst taste:

- application UI before decorative marketing by default;
- production-ready, responsive, accessible components;
- utility-first styling directly in markup;
- neutral, operational visual hierarchy;
- complete interaction states;
- small, composable component APIs;
- conventional primitives before novel layout inventions.
