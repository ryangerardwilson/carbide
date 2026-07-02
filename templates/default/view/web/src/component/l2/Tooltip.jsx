import { useState } from 'react';

export function Tooltip({ children, text }) {
  const [open, setOpen] = useState(false);

  return (
    <span className="relative inline-flex" onBlur={() => setOpen(false)} onFocus={() => setOpen(true)} onMouseEnter={() => setOpen(true)} onMouseLeave={() => setOpen(false)}>
      {children}
      {open ? (
        <span className="absolute bottom-full left-1/2 z-30 mb-2 w-max max-w-56 -translate-x-1/2 rounded-md bg-[#16211b] px-2 py-1 text-xs font-bold text-white shadow-lg">
          {text}
        </span>
      ) : null}
    </span>
  );
}
