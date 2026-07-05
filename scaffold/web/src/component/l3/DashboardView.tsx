import { Metric, Muted, Panel, ui } from '../l1';
import { DashboardLayout } from '../l2';
import { cx } from '../../lib/cx';
import type { NavItem, ResolvedTheme, ThemeMode, User } from '../../lib/types';

const screenClassLayers = {
  workspace: {
    l1: 'grid',
    l2: 'gap-3',
    l3: ''
  },
  statusGrid: {
    l1: 'grid overflow-hidden',
    l2: 'gap-px rounded-lg border md:grid-cols-3',
    l3: cx(ui.border, ui.gridLines)
  },
  statusCell: {
    l1: '',
    l2: 'p-3',
    l3: ui.surface
  },
  sessionLabel: {
    l1: '',
    l2: 'm-0 text-xs font-bold uppercase tracking-normal',
    l3: ui.accent
  },
  sessionTitle: {
    l1: '',
    l2: 'm-0 mt-1 text-base/6 font-semibold',
    l3: ui.text
  }
};

const dashboardNav: NavItem[] = [
  {
    description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
    eyebrow: 'Lorem',
    label: 'Ipsum',
    value: 'workspace'
  }
];

const workspaceMetrics = [
  { label: 'Lorem', value: 'Ipsum dolor' },
  { label: 'Sit', value: 'Amet' },
  { label: 'Consectetur', value: 'Adipiscing' }
];

function WorkspaceOverview() {
  return (
    <div className={cx(screenClassLayers.workspace.l1, screenClassLayers.workspace.l2, screenClassLayers.workspace.l3)}>
      <section
        className={cx(screenClassLayers.statusGrid.l1, screenClassLayers.statusGrid.l2, screenClassLayers.statusGrid.l3)}
        aria-label="Application status"
      >
        {workspaceMetrics.map((metric) => (
          <div className={cx(screenClassLayers.statusCell.l1, screenClassLayers.statusCell.l2, screenClassLayers.statusCell.l3)} key={metric.label}>
            <Metric label={metric.label} value={metric.value} />
          </div>
        ))}
      </section>

      <Panel className="max-w-2xl">
        <p className={cx(screenClassLayers.sessionLabel.l1, screenClassLayers.sessionLabel.l2, screenClassLayers.sessionLabel.l3)}>Dolor</p>
        <h2 className={cx(screenClassLayers.sessionTitle.l1, screenClassLayers.sessionTitle.l2, screenClassLayers.sessionTitle.l3)}>Lorem ipsum dolor sit amet.</h2>
        <Muted className="mt-2 text-xs/5">
          Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
        </Muted>
      </Panel>
    </div>
  );
}

interface DashboardViewProps {
  appName: string;
  busy: boolean;
  onLogout: () => void | Promise<void>;
  onThemeMode: (mode: ThemeMode) => void;
  resolvedTheme: ResolvedTheme;
  themeMode: ThemeMode;
  user: User;
}

export function DashboardView({ appName, busy, onLogout, onThemeMode, resolvedTheme, themeMode, user }: DashboardViewProps) {
  return (
    <DashboardLayout
      activeItem="workspace"
      appName={appName}
      busy={busy}
      navItems={dashboardNav}
      onLogout={onLogout}
      onThemeMode={onThemeMode}
      resolvedTheme={resolvedTheme}
      themeMode={themeMode}
      userEmail={user.email}
    >
      <WorkspaceOverview />
    </DashboardLayout>
  );
}
