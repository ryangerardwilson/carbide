import { useState } from 'react';
import { cx } from '../utils.js';

export function Tabs({ defaultValue, tabs = [] }) {
  const [active, setActive] = useState(defaultValue || tabs[0]?.value);
  const current = tabs.find((tab) => tab.value === active) || tabs[0];

  return (
    <section className="grid gap-4">
      <div className="flex flex-wrap gap-2 border-b border-emerald-950/10" role="tablist">
        {tabs.map((tab) => (
          <button
            aria-selected={tab.value === current?.value}
            className={cx(
              'min-h-10 border-b-2 px-3 text-sm font-bold transition',
              tab.value === current?.value
                ? 'border-teal-700 text-teal-800'
                : 'border-transparent text-[#66786e] hover:text-[#16211b]'
            )}
            key={tab.value}
            role="tab"
            type="button"
            onClick={() => setActive(tab.value)}
          >
            {tab.label}
          </button>
        ))}
      </div>
      <div role="tabpanel">{current?.content}</div>
    </section>
  );
}
