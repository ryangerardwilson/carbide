import { CodeText, Metric, Muted, Panel, ui } from '../l1/index.js';
import { DashboardLayout, Tabs } from '../l2/index.js';
import { cx } from '../utils.js';
import { ComponentLibraryView } from './ComponentLibraryView.jsx';

export function DashboardView({ appName, busy, onLogout, user }) {
  return (
    <DashboardLayout appName={appName} busy={busy} onLogout={onLogout} userEmail={user.email}>
      <Tabs
        tabs={[
          {
            label: 'Workspace',
            value: 'workspace',
            content: (
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
            )
          },
          {
            label: 'Components',
            value: 'components',
            content: <ComponentLibraryView />
          }
        ]}
      />
    </DashboardLayout>
  );
}
