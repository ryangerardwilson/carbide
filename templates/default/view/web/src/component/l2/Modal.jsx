import { useEffect } from 'react';
import { Button, Heading, Muted, Panel } from '../l1/index.js';

export function Modal({ children, description = '', onClose, open, title = 'Dialog' }) {
  useEffect(() => {
    if (!open) {
      return undefined;
    }
    const onKeyDown = (event) => {
      if (event.key === 'Escape') {
        onClose?.();
      }
    };
    window.addEventListener('keydown', onKeyDown);
    return () => window.removeEventListener('keydown', onKeyDown);
  }, [onClose, open]);

  if (!open) {
    return null;
  }

  return (
    <div className="fixed inset-0 z-40 grid place-items-center bg-[#10201c]/45 px-4 py-8" role="presentation">
      <Panel aria-modal="true" className="w-full max-w-lg p-6" role="dialog">
        <div className="mb-5 flex items-start justify-between gap-5">
          <div>
            <Heading level={3}>{title}</Heading>
            {description ? <Muted className="mt-1">{description}</Muted> : null}
          </div>
          <Button aria-label="Close dialog" onClick={onClose} size="sm" variant="ghost">
            Close
          </Button>
        </div>
        {children}
      </Panel>
    </div>
  );
}

export function Slideover({ children, onClose, open, title = 'Panel' }) {
  if (!open) {
    return null;
  }

  return (
    <div className="fixed inset-0 z-40 bg-[#10201c]/40" role="presentation">
      <aside className="ml-auto flex h-full w-full max-w-md flex-col bg-white p-6 shadow-xl" role="dialog" aria-modal="true">
        <div className="mb-5 flex items-start justify-between gap-5">
          <Heading level={3}>{title}</Heading>
          <Button aria-label="Close panel" onClick={onClose} size="sm" variant="ghost">
            Close
          </Button>
        </div>
        {children}
      </aside>
    </div>
  );
}

export const Dialog = Modal;
