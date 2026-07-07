import React from 'react';
import { createRoot } from 'react-dom/client';
import './tailwind.css';

function DocsRuntime() {
  return (
    <main className="min-h-screen bg-amber-50 text-neutral-950 [scrollbar-width:thin] [scrollbar-color:rgb(217_119_6)_transparent] hover:[scrollbar-color:rgb(180_83_9)_transparent] dark:bg-neutral-950 dark:text-amber-50 dark:[scrollbar-color:rgb(250_204_21)_transparent] dark:hover:[scrollbar-color:rgb(253_224_71)_transparent]">
      <section className="mx-auto grid min-h-screen max-w-3xl place-items-center px-4 py-10 text-center">
        <div className="grid gap-3">
          <p className="text-xs font-semibold uppercase tracking-normal text-amber-700 dark:text-yellow-300">
            Carbide Docs
          </p>
          <h1 className="text-2xl/8 font-semibold text-neutral-950 dark:text-amber-50 sm:text-3xl/9">
            Documentation shell ready.
          </h1>
          <p className="text-sm/6 text-amber-950/78 dark:text-amber-200/74">
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
