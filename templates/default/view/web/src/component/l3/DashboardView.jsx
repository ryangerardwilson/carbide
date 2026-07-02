import { CodeText, Metric, Muted, Panel } from '../l1/index.js';
import { DashboardLayout, Tabs } from '../l2/index.js';
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
                  className="grid gap-px overflow-hidden rounded-lg border border-emerald-950/10 bg-emerald-950/10 md:grid-cols-3"
                  aria-label="Application status"
                >
                  {[
                    ['Frontend', 'React + Bun container'],
                    ['Backend', 'Go API container'],
                    ['Database', 'Postgres container']
                  ].map(([label, value]) => (
                    <div className="bg-white p-6" key={label}>
                      <Metric label={label} value={value} />
                    </div>
                  ))}
                </section>

                <Panel className="max-w-3xl">
                  <p className="m-0 text-xs font-extrabold uppercase tracking-normal text-teal-700">Session</p>
                  <h2 className="m-0 mt-2 text-3xl leading-tight text-[#16211b]">Logged in as {user.email}</h2>
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
