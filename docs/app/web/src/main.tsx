import React from 'react';
import { createRoot } from 'react-dom/client';
import './tailwind.css';

function DocsRuntime() {
  return (
    <main className="min-h-screen bg-carbide-page text-carbide-text">
      <section className="mx-auto grid min-h-screen max-w-3xl place-items-center px-4 py-10 text-center">
        <div className="grid gap-3">
          <p className="text-xs font-semibold uppercase tracking-normal text-carbide-muted">
            Carbide Docs
          </p>
          <h1 className="text-2xl/8 font-semibold text-carbide-text sm:text-3xl/9">
            Documentation shell ready.
          </h1>
          <p className="text-sm/6 text-carbide-muted">
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
