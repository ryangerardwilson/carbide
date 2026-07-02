import { cx } from '../utils.js';

export function Toggle({ checked, label, onChange }) {
  return (
    <button
      aria-pressed={checked}
      className="inline-flex min-h-10 items-center gap-3 rounded-md text-left font-bold text-[#16211b]"
      type="button"
      onClick={() => onChange?.(!checked)}
    >
      <span
        className={cx(
          'relative h-6 w-11 rounded-full transition',
          checked ? 'bg-teal-700' : 'bg-emerald-950/20'
        )}
        aria-hidden="true"
      >
        <span
          className={cx(
            'absolute top-1 h-4 w-4 rounded-full bg-white transition',
            checked ? 'left-6' : 'left-1'
          )}
        />
      </span>
      {label}
    </button>
  );
}

export const Switch = Toggle;
