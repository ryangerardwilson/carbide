import { cx } from '../utils.js';
import { ui } from './tokens.js';

export function Field({ children, error = '', hint = '', label, className = '' }) {
  return (
    <label className={cx('grid gap-2 font-bold', ui.text, className)}>
      <span>{label}</span>
      {children}
      {hint && !error ? <span className={cx('text-sm font-normal', ui.subtle)}>{hint}</span> : null}
      {error ? <span className={cx('text-sm font-bold', ui.errorText)}>{error}</span> : null}
    </label>
  );
}

export function TextInput({ className = '', ...props }) {
  return (
    <input
      className={cx(
        'min-h-12 w-full rounded-md border px-3 py-2 outline-none transition',
        ui.input,
        className
      )}
      {...props}
    />
  );
}
