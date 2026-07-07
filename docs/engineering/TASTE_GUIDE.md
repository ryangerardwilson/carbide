# Taste Guide

This file describes Carbide's current starter taste.

Taste is allowed to evolve. It influences scaffold shape, starter UX, package
choices, runtime pins, UI organization, docs examples, and helper defaults.
Taste is not the same thing as law.

## Current Taste Areas

Agents should cite these clauses as `Taste 1` through `Taste 6`.

### Taste 1. Starter Stack

React, Bun, Tailwind, Go, and Postgres are the current starter stack.

### Taste 2. Runtime Pins

Current pinned runtime images and package versions live in scaffolded files.

### Taste 3. Starter Product Surface

The generated starter includes auth, registration, dashboard, and theme
behavior.

### Taste 4. Frontend Organization

Current frontend component organization and Tailwind usage style include
`web/src/component/l1`, `l2`, and `l3` plus L1/L2/L3 class ownership inside
reusable components.

### Taste 5. Docs And Examples

Current docs wording and example flows follow the starter.

### Taste 6. CLI Presentation

Current CLI presentation plus audit and reporting format are part of starter
taste.

## 2. Enforcement Rule

Generated apps may diverge from this taste.

`carbide health` should not fail only because an app no longer matches current
starter taste. Taste is for scaffolding and for audits driven by Codex with
user intent.

## 3. Current Frontend Taste

The current starter frontend teaches two related conventions:

### 3.1 Directory Tiers

- `component/l1`: primitives and Tailwind utility tokens;
- `component/l2`: reusable composed patterns and layouts;
- `component/l3`: product screens and domain-specific sections.

### 3.2 Class Ownership

- `l1`: structure and layout;
- `l2`: geometry, spacing, borders, radii, and type scale;
- `l3`: theme, color, state, motion, and interaction.

### 3.3 Composition Rule

When a component is reusable or variant-heavy, prefer explicit class-layer
maps and `cx()` composition over long unreadable inline class strings.

### 3.4 Tailwind Plus / Catalyst Component Taste

Component design should also stay close to Tailwind Plus / Catalyst taste:

- application UI before decorative marketing by default;
- production-ready, responsive, accessible components;
- utility-first styling directly in markup;
- neutral, operational visual hierarchy;
- complete interaction states;
- small, composable component APIs;
- conventional primitives before novel layout inventions.
