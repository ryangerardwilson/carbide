import { CodeText, Metric, Muted, Panel, ui } from '../l1/index.js';
import { DashboardLayout } from '../l2/index.js';
import { cx } from '../../lib/cx.js';

const screenClassLayers = {
  workspace: {
    l1: 'grid',
    l2: 'gap-6',
    l3: ''
  },
  statusGrid: {
    l1: 'grid overflow-hidden',
    l2: 'gap-px rounded-lg border md:grid-cols-3',
    l3: cx(ui.border, ui.gridLines)
  },
  statusCell: {
    l1: '',
    l2: 'p-6',
    l3: ui.surface
  },
  sessionLabel: {
    l1: '',
    l2: 'm-0 text-xs font-extrabold uppercase tracking-normal',
    l3: ui.accent
  },
  sessionTitle: {
    l1: '',
    l2: 'm-0 mt-2 text-3xl leading-tight',
    l3: ui.text
  }
};

const dashboardNav = [
  {
    description: 'Runtime health, same-origin API flow, and the active session.',
    eyebrow: 'Local app',
    label: 'Workspace',
    value: 'workspace'
  }
];

function WorkspaceOverview({ user }) {
  return (
    <div className={cx(screenClassLayers.workspace.l1, screenClassLayers.workspace.l2, screenClassLayers.workspace.l3)}>
      <section
        className={cx(screenClassLayers.statusGrid.l1, screenClassLayers.statusGrid.l2, screenClassLayers.statusGrid.l3)}
        aria-label="Application status"
      >
        {[
          ['Web', 'React + Bun container'],
          ['API', 'Go API container'],
          ['Database', 'Postgres db container']
        ].map(([label, value]) => (
          <div className={cx(screenClassLayers.statusCell.l1, screenClassLayers.statusCell.l2, screenClassLayers.statusCell.l3)} key={label}>
            <Metric label={label} value={value} />
          </div>
        ))}
      </section>

      <Panel className="max-w-3xl">
        <p className={cx(screenClassLayers.sessionLabel.l1, screenClassLayers.sessionLabel.l2, screenClassLayers.sessionLabel.l3)}>Session</p>
        <h2 className={cx(screenClassLayers.sessionTitle.l1, screenClassLayers.sessionTitle.l2, screenClassLayers.sessionTitle.l3)}>Logged in as {user.email}</h2>
        <Muted className="mt-4">
          The browser talks to <CodeText>/api</CodeText> on the same origin. Bun proxies those
          requests to Go, and Go persists the session in Postgres.
        </Muted>
      </Panel>
    </div>
  );
}

export function DashboardView({ appName, busy, onLogout, user }) {
  return (
    <DashboardLayout
      activeItem="workspace"
      appName={appName}
      busy={busy}
      navItems={dashboardNav}
      onLogout={onLogout}
      userEmail={user.email}
    >
      <WorkspaceOverview user={user} />
    </DashboardLayout>
  );
}
