import { useState } from 'react';
import { CodeText, Metric, Muted, Panel, ui } from '../l1/index.js';
import { DashboardLayout } from '../l2/index.js';
import { cx } from '../utils.js';
import { ComponentLibraryView } from './ComponentLibraryView.jsx';

const dashboardNav = [
  {
    description: 'Runtime health, same-origin API flow, and the active session.',
    eyebrow: 'Local app',
    label: 'Workspace',
    value: 'workspace'
  },
  {
    description: 'Generated L1 primitives, L2 patterns, and L3 starter surfaces.',
    eyebrow: 'React library',
    label: 'Components',
    value: 'components'
  }
];

function WorkspaceOverview({ user }) {
  return (
    <div className="grid gap-6">
      <section
        className={cx('cb-grid-lines grid gap-px overflow-hidden rounded-lg border md:grid-cols-3', ui.border)}
        aria-label="Application status"
      >
        {[
          ['Frontend', 'React + Bun container'],
          ['Backend', 'Go API container'],
          ['Database', 'Postgres container']
        ].map(([label, value]) => (
          <div className={cx('p-6', ui.surface)} key={label}>
            <Metric label={label} value={value} />
          </div>
        ))}
      </section>

      <Panel className="max-w-3xl">
        <p className={cx('m-0 text-xs font-extrabold uppercase tracking-normal', ui.accent)}>Session</p>
        <h2 className={cx('m-0 mt-2 text-3xl leading-tight', ui.text)}>Logged in as {user.email}</h2>
        <Muted className="mt-4">
          The browser talks to <CodeText>/api</CodeText> on the same origin. Bun proxies those
          requests to Go, and Go persists the session in Postgres.
        </Muted>
      </Panel>
    </div>
  );
}

export function DashboardView({ appName, busy, onLogout, user }) {
  const [activeSection, setActiveSection] = useState('workspace');

  return (
    <DashboardLayout
      activeItem={activeSection}
      appName={appName}
      busy={busy}
      navItems={dashboardNav}
      onLogout={onLogout}
      onNavItem={setActiveSection}
      userEmail={user.email}
    >
      {activeSection === 'components' ? <ComponentLibraryView /> : <WorkspaceOverview user={user} />}
    </DashboardLayout>
  );
}
