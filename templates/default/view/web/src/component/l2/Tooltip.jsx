import { useState } from 'react';
import { cx } from '../utils.js';

export function Tooltip({ children, text }) {
  const [open, setOpen] = useState(false);

  return (
    <span className="relative inline-flex" onBlur={() => setOpen(false)} onFocus={() => setOpen(true)} onMouseEnter={() => setOpen(true)} onMouseLeave={() => setOpen(false)}>
      {children}
      {open ? (
        <span className={cx('cb-tooltip cb-shadow-elevated absolute bottom-full left-1/2 z-30 mb-2 w-max max-w-56 -translate-x-1/2 rounded-md px-2 py-1 text-xs font-bold')}>
          {text}
        </span>
      ) : null}
    </span>
  );
}
