import { Button } from '../l1/Button.jsx';
import { Divider } from '../l1/Surface.jsx';
import { Eyebrow, Heading, Muted } from '../l1/Text.jsx';
import { ui } from '../l1/tokens.js';
import { cx } from '../../lib/cx.js';

const landingClassLayers = {
  shell: {
    l1: 'grid',
    l2: 'min-h-svh lg:grid-cols-2',
    l3: ui.page
  },
  hero: {
    l1: 'grid content-end',
    l2: 'min-h-96 px-8 py-10 sm:px-12 lg:min-h-svh lg:px-16 lg:py-24 xl:px-24',
    l3: ui.hero
  },
  eyebrow: {
    l1: '',
    l2: 'mb-3 text-xs font-extrabold uppercase tracking-normal',
    l3: ui.heroMuted
  },
  title: {
    l1: '',
    l2: 'm-0 max-w-4xl text-5xl leading-none sm:text-6xl lg:text-7xl',
    l3: ''
  },
  copy: {
    l1: '',
    l2: 'mt-5 max-w-2xl text-lg',
    l3: ui.heroMuted
  }
};

const dashboardClassLayers = {
  shell: {
    l1: '',
    l2: 'min-h-svh',
    l3: ui.page
  },
  grid: {
    l1: 'grid',
    l2: 'min-h-svh lg:grid-cols-[280px_minmax(0,1fr)]',
    l3: ''
  },
  sidebar: {
    l1: 'flex min-w-0 flex-col lg:sticky',
    l2: 'border-b px-5 py-5 lg:top-0 lg:h-svh lg:border-b-0 lg:border-r',
    l3: cx(ui.border, ui.surface)
  },
  nav: {
    l1: 'flex overflow-x-auto lg:grid lg:overflow-visible',
    l2: 'mt-6 gap-2 pb-1 lg:pb-0',
    l3: ''
  },
  navButton: {
    l1: 'shrink-0',
    l2: 'min-h-11 rounded-md border px-3 text-left text-sm font-bold lg:w-full',
    l3: cx(ui.focus, 'transition')
  },
  dividerWrap: {
    l1: 'hidden lg:block',
    l2: 'mt-6',
    l3: ''
  },
  footer: {
    l1: '',
    l2: 'mt-5 lg:mt-auto',
    l3: ''
  },
  content: {
    l1: 'min-w-0',
    l2: 'px-6 py-8 sm:px-10 lg:py-12',
    l3: ''
  },
  header: {
    l1: '',
    l2: 'mb-8 border-b pb-7',
    l3: ui.border
  }
};

export function LandingPageLayout({ appName, children, mode }) {
  const isRegister = mode === 'register';

  return (
    <main className={cx(landingClassLayers.shell.l1, landingClassLayers.shell.l2, landingClassLayers.shell.l3)}>
      <section className={cx(landingClassLayers.hero.l1, landingClassLayers.hero.l2, landingClassLayers.hero.l3)}>
        <p className={cx(landingClassLayers.eyebrow.l1, landingClassLayers.eyebrow.l2, landingClassLayers.eyebrow.l3)}>{appName}</p>
        <h1 className={cx(landingClassLayers.title.l1, landingClassLayers.title.l2, landingClassLayers.title.l3)}>
          {isRegister ? 'Create the first account.' : 'Log in to the workspace.'}
        </h1>
        <p className={cx(landingClassLayers.copy.l1, landingClassLayers.copy.l2, landingClassLayers.copy.l3)}>
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
    <main className={cx(dashboardClassLayers.shell.l1, dashboardClassLayers.shell.l2, dashboardClassLayers.shell.l3)}>
      <div className={cx(dashboardClassLayers.grid.l1, dashboardClassLayers.grid.l2, dashboardClassLayers.grid.l3)}>
        <aside
          className={cx(
            dashboardClassLayers.sidebar.l1,
            dashboardClassLayers.sidebar.l2,
            dashboardClassLayers.sidebar.l3
          )}
        >
          <div>
            <Eyebrow>Bun + Go + Postgres</Eyebrow>
            <Heading className="mt-2" level={2}>
              {appName}
            </Heading>
            <Muted className="mt-2 text-sm">Signed in as {userEmail}</Muted>
          </div>

          {navItems.length ? (
            <nav className={cx(dashboardClassLayers.nav.l1, dashboardClassLayers.nav.l2, dashboardClassLayers.nav.l3)} aria-label="Dashboard">
              {navItems.map((item) => {
                const active = item.value === activeItem;

                return (
                  <button
                    aria-current={active ? 'page' : undefined}
                    className={cx(
                      dashboardClassLayers.navButton.l1,
                      dashboardClassLayers.navButton.l2,
                      active ? ui.selectionActive : ui.selection,
                      dashboardClassLayers.navButton.l3
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

          <div className={cx(dashboardClassLayers.dividerWrap.l1, dashboardClassLayers.dividerWrap.l2, dashboardClassLayers.dividerWrap.l3)}>
            <Divider />
          </div>

          <div className={cx(dashboardClassLayers.footer.l1, dashboardClassLayers.footer.l2, dashboardClassLayers.footer.l3)}>
            <Button className="w-full" disabled={busy} onClick={onLogout} variant="secondary">
              {busy ? 'Logging out...' : 'Log out'}
            </Button>
          </div>
        </aside>

        <section className={cx(dashboardClassLayers.content.l1, dashboardClassLayers.content.l2, dashboardClassLayers.content.l3)}>
          <header className={cx(dashboardClassLayers.header.l1, dashboardClassLayers.header.l2, dashboardClassLayers.header.l3)}>
            <Eyebrow>{activeNavItem?.eyebrow || 'Dashboard'}</Eyebrow>
            <Heading className="mt-2">{activeNavItem?.label || 'Workspace'}</Heading>
            {activeNavItem?.description ? <Muted className="mt-3 max-w-3xl">{activeNavItem.description}</Muted> : null}
          </header>
          {children}
        </section>
      </div>
    </main>
  );
}
