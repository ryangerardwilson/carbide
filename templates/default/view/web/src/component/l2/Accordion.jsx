import { useState } from 'react';
import { ui } from '../l1/index.js';
import { cx } from '../utils.js';

export function Accordion({ defaultOpen = 0, items = [] }) {
  const [open, setOpen] = useState(defaultOpen);

  return (
    <div className={cx('cb-divide-y rounded-lg border', ui.border, ui.surface)}>
      {items.map((item, index) => (
        <section key={item.title}>
          <button
            aria-expanded={open === index}
            className={cx('flex min-h-12 w-full items-center justify-between gap-4 px-4 text-left font-bold', ui.text)}
            type="button"
            onClick={() => setOpen(open === index ? -1 : index)}
          >
            <span>{item.title}</span>
            <span aria-hidden="true">{open === index ? '-' : '+'}</span>
          </button>
          {open === index ? <div className={cx('px-4 pb-4 text-sm', ui.subtle)}>{item.content}</div> : null}
        </section>
      ))}
    </div>
  );
}

export const Disclosure = Accordion;
