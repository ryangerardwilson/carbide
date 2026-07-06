# Tailwind Component Organization

The generated web app uses Tailwind as the required styling path. Keep UI
changes inside the Tailwind utility system unless a repeated pattern clearly
belongs in `web/src/component/l1/tokens.ts`. Do not create a generated CSS
variable palette for colors that can live as ordinary Tailwind classes.

## File Ownership

- `web/src/component/l1/` owns primitives: buttons, inputs, surfaces, text,
  the theme toggle, and tokens. L1 components must not know app-specific data,
  routes, auth flow, fetches, sessions, or product domain facts.
- `web/src/component/l2/` owns reusable composed patterns: forms, app shells,
  page layouts, and repeated UX structures made from L1 primitives. L2 may own
  presentation interaction state, but not backend behavior or product rules.
- `web/src/component/l3/` owns product screens and domain sections composed
  from L2 patterns and L1 primitives. L3 may know authenticated user data,
  dashboard copy, and route-level product facts.
- `web/src/component/l1/tokens.ts` owns stable names for shared Tailwind
  utility groups, including the starter light/dark visual choices.
- `web/src/component/l1/tokens.ts` also owns the built-in scrollbar utility
  group. Use standards-based `scrollbar-width` and `scrollbar-color` utilities
  with light/dark variants; do not add `::-webkit-scrollbar` or other
  scrollbar pseudo-selector CSS to `web/src/styles.css`.
- `web/src/styles.css` owns only the Tailwind import, TypeScript-aware
  `@source` directives, and the `data-theme` dark variant.
- Do not put global `html`/`body` font-size, width, line-height, layout, color,
  or component styling rules in `web/src/styles.css`; use Tailwind utilities
  and component tokens instead.
- `web/index.html` owns the no-flash theme bootstrap before React loads.
- `web/src/lib/cx.ts` owns the small class-name helper used by components.

## Tailwind Class Layers

Use L1/L2/L3 in two related ways:

- Directory ownership: `l1/` primitives, `l2/` reusable patterns, and `l3/`
  product screens.
- Class ownership inside each component:

- L1 structure: stable layout and element behavior such as `relative`, `flex`,
  `grid`, `overflow-hidden`, and `box-border`.
- L2 geometry: size, spacing, borders, radius, and type scale such as `h-10`,
  `p-4`, `rounded-md`, `border`, `text-sm`, and `font-semibold`.
- L3 theme and state: colors, shadows, interaction, animation, focus, disabled,
  and dark-mode behavior such as `bg-*`, `text-*`, `shadow-*`, `hover:*`,
  `focus-visible:*`, `disabled:*`, and `dark:*`.

Refactor when a class string becomes hard to scan:

- Keep L1 structural classes in the stable base string.
- Move L2 choices into `size`, `density`, or `shape` variants.
- Move L3 choices into `variant`, `intent`, `tone`, or state variants.
- Accept `className` for caller placement and rare overrides, not as the normal
  variant API.
- Use the existing `cx()` helper for small conditional class lists.
- Add heavier helpers such as `clsx`, `tailwind-merge`, or
  `class-variance-authority` only when repeated variants justify the dependency.

## Class Ordering

Group class names by purpose so diffs remain readable:

1. Layout and box model.
2. Flex/grid children and alignment.
3. Sizing.
4. Spacing.
5. Typography.
6. Background, color, fill, and stroke.
7. Border, radius, shadow, ring, and outline.
8. Transform, transition, and animation.
9. Interactivity without variants.
10. Variant modifiers such as `hover:`, `focus-visible:`, `disabled:`,
    `aria-*`, `data-*`, `dark:`, and responsive prefixes.

Use mobile-first classes. Write the base state for the smallest viewport, then
add `sm:`, `md:`, `lg:`, `xl:`, and `2xl:` variants only as needed.

## Accessibility Rules

- Every interactive element needs visible focus styling.
- Prefer semantic HTML before ARIA roles.
- Pair inputs with visible labels.
- Use `aria-hidden="true"` for decorative icons and `sr-only` text for
  icon-only controls.
- Disabled states must be obvious visually and use the correct disabled
  attribute.

## CSS And Theme Rules

- Prefer Tailwind utilities for one-off layout and styling.
- Keep generated colors, focus rings, borders, and light/dark variants in
  `web/src/component/l1/tokens.ts`, using Tailwind utility classes.
- Do not add generated `--carbide-*` color variables or a generated `@theme`
  palette in `web/src/styles.css`.
- Keep theme choice as browser-local React state plus `localStorage`; do not
  send theme preference to the API or database by default.
- Use `ThemeToggle.tsx` for the built-in `light`, `dark`, and `system`
  controls instead of creating one-off theme switches in screens.
- Keep repeated component utility groups in `tokens.ts`.
- Do not add a parallel `theme.css` file or custom `cb-*` component classes.
- Use custom CSS only for behavior that utilities cannot express cleanly.
- Avoid arbitrary values unless the exact value is required by the design.
- `carbide doctor` must fail when `web/src/styles.css` grows into component
  CSS. Custom class selectors, ID selectors, `@apply`, `@layer`, keyframes, and
  media/container rules belong in Tailwind utility classes, component class
  layers, or a deliberate component abstraction instead.
