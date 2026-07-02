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

export function DashboardLayout({
  activeItem = '',
  appName,
  busy,
  children,
  navItems = [],
  onLogout,
  onNavItem,
  userEmail
}) {
  const activeNavItem = navItems.find((item) => item.value === activeItem) || navItems[0];

  return (
    <main className={cx('min-h-svh', ui.page)}>
      <div className="grid min-h-svh lg:grid-cols-[280px_minmax(0,1fr)]">
        <aside
          className={cx(
            'flex min-w-0 flex-col border-b px-5 py-5 lg:sticky lg:top-0 lg:h-svh lg:border-b-0 lg:border-r',
            ui.border,
            ui.surface
          )}
        >
          <div>
            <Eyebrow>Bun + Go + Postgres</Eyebrow>
            <Heading className="mt-2 text-3xl" level={2}>
              {appName}
            </Heading>
            <Muted className="mt-2 text-sm">Signed in as {userEmail}</Muted>
          </div>

          {navItems.length ? (
            <nav className="mt-6 flex gap-2 overflow-x-auto pb-1 lg:grid lg:overflow-visible lg:pb-0" aria-label="Dashboard">
              {navItems.map((item) => {
                const active = item.value === activeItem;

                return (
                  <button
                    aria-current={active ? 'page' : undefined}
                    className={cx(
                      'min-h-11 shrink-0 rounded-md border px-3 text-left text-sm font-bold transition lg:w-full',
                      active ? 'cb-selection-active' : 'cb-selection'
                    )}
                    key={item.value}
                    type="button"
                    onClick={() => onNavItem?.(item.value)}
                  >
                    {item.label}
                  </button>
                );
              })}
            </nav>
          ) : null}

          <div className="mt-6 hidden lg:block">
            <div className={cx('h-px w-full', ui.divider)} />
          </div>

          <div className="mt-5 lg:mt-auto">
            <Button className="w-full" disabled={busy} onClick={onLogout} variant="secondary">
              {busy ? 'Logging out...' : 'Log out'}
            </Button>
          </div>
        </aside>

        <section className="min-w-0 px-6 py-8 sm:px-10 lg:py-12">
          <header className={cx('mb-8 border-b pb-7', ui.border)}>
            <Eyebrow>{activeNavItem?.eyebrow || 'Dashboard'}</Eyebrow>
            <Heading className="mt-2 text-4xl sm:text-5xl">{activeNavItem?.label || 'Workspace'}</Heading>
            {activeNavItem?.description ? <Muted className="mt-3 max-w-3xl">{activeNavItem.description}</Muted> : null}
          </header>
          {children}
        </section>
      </div>
    </main>
  );
}
