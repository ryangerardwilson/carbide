import { useState } from 'react';
import { Button, Panel } from '../l1/index.js';

export function Popover({ children, label = 'Open popover' }) {
  const [open, setOpen] = useState(false);

  return (
    <div className="relative inline-flex">
      <Button aria-expanded={open} onClick={() => setOpen((value) => !value)} size="sm" variant="secondary">
        {label}
      </Button>
      {open ? (
        <Panel as="div" className="absolute left-0 top-full z-30 mt-2 w-72 p-4 text-left">
          {children}
        </Panel>
      ) : null}
    </div>
  );
}
