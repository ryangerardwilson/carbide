import type { ReactNode } from 'react';
import { Button } from '../l1/Button';
import { Divider } from '../l1/Surface';
import { ThemeToggle } from '../l1/ThemeToggle';
import { Eyebrow, Heading, Muted } from '../l1/Text';
import { ui } from '../l1/tokens';
import { cx } from '../../lib/cx';
import type { NavItem, ResolvedTheme, ThemeMode } from '../../lib/types';

const landingClassLayers = {
  shell: {
    l1: 'relative grid',
    l2: 'min-h-svh lg:grid-cols-2',
    l3: ui.page
  },
  hero: {
    l1: 'grid content-between',
    l2: 'min-h-64 px-4 py-5 sm:px-6 lg:min-h-svh lg:px-8 lg:py-10 xl:px-10',
    l3: ui.hero
  },
  heroBody: {
    l1: '',
    l2: '',
    l3: ''
  },
  eyebrow: {
    l1: '',
    l2: 'mb-2 text-xs font-bold uppercase tracking-normal',
    l3: ui.heroMuted
  },
  title: {
    l1: '',
    l2: 'm-0 max-w-2xl text-2xl/8 sm:text-3xl/9 lg:text-4xl/10',
    l3: ''
  },
  copy: {
    l1: '',
    l2: 'mt-2 max-w-lg text-xs/5 sm:text-sm/6',
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
    l2: 'min-h-svh lg:grid-cols-[216px_minmax(0,1fr)]',
    l3: ''
  },
  sidebar: {
    l1: 'flex min-w-0 flex-col lg:sticky',
    l2: 'border-b px-3 py-3 lg:top-0 lg:h-svh lg:border-b-0 lg:border-r',
    l3: cx(ui.border, ui.surface)
  },
  nav: {
    l1: 'flex overflow-x-auto lg:grid lg:overflow-visible',
    l2: 'mt-3 gap-1 pb-1 lg:pb-0',
    l3: ''
  },
  navButton: {
    l1: 'shrink-0',
    l2: 'min-h-8 rounded-md border px-2 text-left text-xs font-semibold lg:w-full',
    l3: cx(ui.focus, 'transition')
  },
  dividerWrap: {
    l1: 'hidden lg:block',
    l2: 'mt-3',
    l3: ''
  },
  footer: {
    l1: '',
    l2: 'mt-3 lg:mt-auto',
    l3: ''
  },
  content: {
    l1: 'min-w-0',
    l2: 'px-3 py-4 sm:px-5 lg:py-5',
    l3: ''
  },
  header: {
    l1: '',
    l2: 'mb-4 border-b pb-3',
    l3: ui.border
  }
};

interface ThemeControlProps {
  onThemeMode: (mode: ThemeMode) => void;
  resolvedTheme: ResolvedTheme;
  themeMode: ThemeMode;
}

interface LandingPageLayoutProps extends ThemeControlProps {
  appName: string;
  children: ReactNode;
}

interface DashboardLayoutProps extends ThemeControlProps {
  activeItem?: string;
  appName: string;
  busy: boolean;
  children: ReactNode;
  navItems?: NavItem[];
  onLogout: () => void | Promise<void>;
  onNavItem?: (value: string) => void;
  userEmail: string;
}

export function LandingPageLayout({ appName, children, onThemeMode, resolvedTheme, themeMode }: LandingPageLayoutProps) {
  return (
    <main className={cx(landingClassLayers.shell.l1, landingClassLayers.shell.l2, landingClassLayers.shell.l3)}>
      <ThemeToggle
        className="absolute right-4 top-4 z-10 sm:right-6 lg:right-8 lg:top-6"
        mode={themeMode}
        onMode={onThemeMode}
        resolved={resolvedTheme}
      />
      <section className={cx(landingClassLayers.hero.l1, landingClassLayers.hero.l2, landingClassLayers.hero.l3)}>
        <div className={cx(landingClassLayers.heroBody.l1, landingClassLayers.heroBody.l2, landingClassLayers.heroBody.l3)}>
          <p className={cx(landingClassLayers.eyebrow.l1, landingClassLayers.eyebrow.l2, landingClassLayers.eyebrow.l3)}>{appName}</p>
          <h1 className={cx(landingClassLayers.title.l1, landingClassLayers.title.l2, landingClassLayers.title.l3)}>
            Lorem ipsum dolor sit amet.
          </h1>
          <p className={cx(landingClassLayers.copy.l1, landingClassLayers.copy.l2, landingClassLayers.copy.l3)}>
            Consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore.
          </p>
        </div>
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
  onThemeMode,
  resolvedTheme,
  themeMode,
  userEmail
}: DashboardLayoutProps) {
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
            <Eyebrow>Workspace</Eyebrow>
            <Heading className="mt-1.5" level={2}>
              {appName}
            </Heading>
            <Muted className="mt-1 text-xs/5">Signed in as {userEmail}</Muted>
            <ThemeToggle className="mt-3" mode={themeMode} onMode={onThemeMode} resolved={resolvedTheme} />
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
            <Heading className="mt-1.5">{activeNavItem?.label || 'Workspace'}</Heading>
            {activeNavItem?.description ? <Muted className="mt-1.5 max-w-xl text-xs/5">{activeNavItem.description}</Muted> : null}
          </header>
          {children}
        </section>
      </div>
    </main>
  );
}
