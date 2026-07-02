import { ui } from '../l1/index.js';
import { cx } from '../utils.js';

export function Toggle({ checked, label, onChange }) {
  return (
    <button
      aria-pressed={checked}
      className={cx('inline-flex min-h-10 items-center gap-3 rounded-md text-left font-bold', ui.text)}
      type="button"
      onClick={() => onChange?.(!checked)}
    >
      <span
        className={cx(
          'relative h-6 w-11 rounded-full transition',
          checked ? 'cb-switch-track-on' : 'cb-switch-track-off'
        )}
        aria-hidden="true"
      >
        <span
          className={cx(
            'absolute top-1 h-4 w-4 rounded-full transition',
            ui.surface,
            checked ? 'left-6' : 'left-1'
          )}
        />
      </span>
      {label}
    </button>
  );
}

export const Switch = Toggle;
