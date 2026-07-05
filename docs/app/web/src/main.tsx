import React from 'react';
import { createRoot } from 'react-dom/client';
import './tailwind.css';

function DocsRuntime() {
  return (
    <main className="min-h-screen bg-white text-neutral-950 dark:bg-black dark:text-neutral-50">
      <section className="mx-auto grid min-h-screen max-w-3xl place-items-center px-4 py-10 text-center">
        <div className="grid gap-3">
          <p className="text-xs font-semibold uppercase tracking-normal text-neutral-600 dark:text-neutral-400">
            Carbide Docs
          </p>
          <h1 className="text-2xl/8 font-semibold text-neutral-950 dark:text-neutral-50 sm:text-3xl/9">
            Documentation shell ready.
          </h1>
          <p className="text-sm/6 text-neutral-600 dark:text-neutral-400">
            The production docs routes are served from the checked-in static documentation site.
          </p>
        </div>
      </section>
    </main>
  );
}

const rootElement = document.getElementById('root');
if (rootElement) {
  createRoot(rootElement).render(<DocsRuntime />);
}
