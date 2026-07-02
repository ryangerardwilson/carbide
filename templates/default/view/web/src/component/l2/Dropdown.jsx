import { useEffect, useRef, useState } from 'react';
import { Button, Panel } from '../l1/index.js';

export function Dropdown({ align = 'left', items = [], label = 'Open menu', onSelect }) {
  const [open, setOpen] = useState(false);
  const ref = useRef(null);

  useEffect(() => {
    const onPointerDown = (event) => {
      if (ref.current && !ref.current.contains(event.target)) {
        setOpen(false);
      }
    };
    window.addEventListener('pointerdown', onPointerDown);
    return () => window.removeEventListener('pointerdown', onPointerDown);
  }, []);

  return (
    <div className="relative inline-flex" ref={ref}>
      <Button aria-expanded={open} aria-haspopup="menu" onClick={() => setOpen((value) => !value)} variant="secondary">
        {label}
      </Button>
      {open ? (
        <Panel
          as="div"
          className={`absolute top-full z-30 mt-2 min-w-56 p-1 ${align === 'right' ? 'right-0' : 'left-0'}`}
          role="menu"
        >
          {items.map((item) => (
            <button
              className="block min-h-10 w-full rounded-md px-3 text-left text-sm text-[#16211b] hover:bg-emerald-50"
              key={item.value || item.label}
              role="menuitem"
              type="button"
              onClick={() => {
                onSelect?.(item);
                setOpen(false);
              }}
            >
              {item.label}
            </button>
          ))}
        </Panel>
      ) : null}
    </div>
  );
}

export const Menu = Dropdown;
