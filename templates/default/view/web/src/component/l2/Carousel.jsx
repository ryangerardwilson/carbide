import { useState } from 'react';
import { Button, ui } from '../l1/index.js';
import { cx } from '../utils.js';

export function Carousel({ items = [], label = 'Carousel' }) {
  const [index, setIndex] = useState(0);
  const item = items[index] || {};

  if (!items.length) {
    return null;
  }

  return (
    <section aria-label={label} className="grid gap-3">
      <div className={cx('min-h-36 rounded-lg border p-5', ui.border, ui.surface)}>
        <p className={cx('m-0 text-xs font-bold uppercase', ui.accent)}>
          {index + 1} / {items.length}
        </p>
        <h3 className={cx('m-0 mt-2 text-xl', ui.text)}>{item.title}</h3>
        {item.detail ? <p className={cx('m-0 mt-2 text-sm', ui.subtle)}>{item.detail}</p> : null}
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
