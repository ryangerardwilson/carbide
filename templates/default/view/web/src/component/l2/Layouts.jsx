import { Button, Eyebrow, Heading, Muted, ui } from '../l1/index.js';
import { cx } from '../utils.js';

export function LandingPageLayout({ appName, children, mode }) {
  const isRegister = mode === 'register';

  return (
    <main className={cx('grid min-h-svh lg:grid-cols-[minmax(0,1fr)_minmax(360px,480px)]', ui.page)}>
      <section className="cb-hero cb-on-hero grid min-h-[42svh] content-end px-8 py-10 sm:px-12 lg:min-h-svh lg:px-[7vw] lg:py-[7vw]">
        <p className="cb-on-hero-muted mb-3 text-xs font-extrabold uppercase tracking-normal">{appName}</p>
        <h1 className="m-0 max-w-4xl text-[clamp(42px,7vw,82px)] leading-none">
          {isRegister ? 'Create the first account.' : 'Log in to the workspace.'}
        </h1>
        <p className="cb-on-hero-muted mt-5 max-w-2xl text-lg">
          React and Tailwind own the browser. Go owns the API. Postgres owns durable state.
        </p>
      </section>
      {children}
    </main>
  );
}

export function DashboardLayout({ appName, busy, children, onLogout, userEmail }) {
  return (
    <main className={cx('mx-auto min-h-svh max-w-7xl px-6 py-8 sm:px-10 lg:py-12', ui.page)}>
      <header className={cx('mb-9 flex flex-col gap-5 border-b pb-7 sm:flex-row sm:items-end sm:justify-between', ui.border)}>
        <div>
          <Eyebrow>Bun frontend + Go API + Postgres</Eyebrow>
          <Heading className="mt-2 text-5xl sm:text-6xl">{appName}</Heading>
          <Muted className="mt-3">Signed in as {userEmail}</Muted>
        </div>
        <Button disabled={busy} onClick={onLogout} variant="secondary">
          {busy ? 'Logging out...' : 'Log out'}
        </Button>
      </header>
      {children}
    </main>
  );
}
