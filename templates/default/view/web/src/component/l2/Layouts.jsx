import { Button, Eyebrow, Heading, Muted } from '../l1/index.js';

export function LandingPageLayout({ appName, children, mode }) {
  const isRegister = mode === 'register';

  return (
    <main className="grid min-h-svh bg-[#f6f8f5] lg:grid-cols-[minmax(0,1fr)_minmax(360px,480px)]">
      <section className="grid min-h-[42svh] content-end bg-[linear-gradient(150deg,#0f766e_0%,#1b3f3a_48%,#16211b_100%)] px-8 py-10 text-white sm:px-12 lg:min-h-svh lg:px-[7vw] lg:py-[7vw]">
        <p className="mb-3 text-xs font-extrabold uppercase tracking-normal text-white/75">{appName}</p>
        <h1 className="m-0 max-w-4xl text-[clamp(42px,7vw,82px)] leading-none">
          {isRegister ? 'Create the first account.' : 'Log in to the workspace.'}
        </h1>
        <p className="mt-5 max-w-2xl text-lg text-white/80">
          React and Tailwind own the browser. Go owns the API. Postgres owns durable state.
        </p>
      </section>
      {children}
    </main>
  );
}

export function DashboardLayout({ appName, busy, children, onLogout, userEmail }) {
  return (
    <main className="mx-auto min-h-svh max-w-7xl px-6 py-8 sm:px-10 lg:py-12">
      <header className="mb-9 flex flex-col gap-5 border-b border-emerald-950/10 pb-7 sm:flex-row sm:items-end sm:justify-between">
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
