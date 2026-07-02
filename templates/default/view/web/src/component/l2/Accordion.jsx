import { useState } from 'react';

export function Accordion({ defaultOpen = 0, items = [] }) {
  const [open, setOpen] = useState(defaultOpen);

  return (
    <div className="divide-y divide-emerald-950/10 rounded-lg border border-emerald-950/10 bg-white">
      {items.map((item, index) => (
        <section key={item.title}>
          <button
            aria-expanded={open === index}
            className="flex min-h-12 w-full items-center justify-between gap-4 px-4 text-left font-bold text-[#16211b]"
            type="button"
            onClick={() => setOpen(open === index ? -1 : index)}
          >
            <span>{item.title}</span>
            <span aria-hidden="true">{open === index ? '-' : '+'}</span>
          </button>
          {open === index ? <div className="px-4 pb-4 text-sm text-[#66786e]">{item.content}</div> : null}
        </section>
      ))}
    </div>
  );
}

export const Disclosure = Accordion;
