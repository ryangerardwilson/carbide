import { useState } from 'react';
import { Button } from '../l1/index.js';

export function Carousel({ items = [], label = 'Carousel' }) {
  const [index, setIndex] = useState(0);
  const item = items[index] || {};

  if (!items.length) {
    return null;
  }

  return (
    <section aria-label={label} className="grid gap-3">
      <div className="min-h-36 rounded-lg border border-emerald-950/10 bg-white p-5">
        <p className="m-0 text-xs font-bold uppercase text-teal-700">
          {index + 1} / {items.length}
        </p>
        <h3 className="m-0 mt-2 text-xl text-[#16211b]">{item.title}</h3>
        {item.detail ? <p className="m-0 mt-2 text-sm text-[#66786e]">{item.detail}</p> : null}
      </div>
      <div className="flex gap-2">
        <Button onClick={() => setIndex((index - 1 + items.length) % items.length)} size="sm" variant="secondary">
          Previous
        </Button>
        <Button onClick={() => setIndex((index + 1) % items.length)} size="sm" variant="secondary">
          Next
        </Button>
      </div>
    </section>
  );
}
