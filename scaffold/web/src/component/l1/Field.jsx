import { cx } from '../../lib/cx.js';
import { ui } from './tokens.js';

const fieldClassLayers = {
  l1: 'grid',
  l2: 'gap-1 text-sm/6 font-semibold',
  l3: ui.text
};

const fieldHintClassLayers = {
  l1: '',
  l2: 'text-xs font-normal',
  l3: ui.subtle
};

const fieldErrorClassLayers = {
  l1: '',
  l2: 'text-xs font-bold',
  l3: ui.errorText
};

const inputClassLayers = {
  l1: 'block w-full',
  l2: 'min-h-8 rounded-md border px-2 py-1 text-sm/6 outline-none',
  l3: cx(ui.input, ui.focus, 'transition')
};

export function Field({ children, error = '', hint = '', label, className = '' }) {
  return (
    <label className={cx(fieldClassLayers.l1, fieldClassLayers.l2, fieldClassLayers.l3, className)}>
      <span>{label}</span>
      {children}
      {hint && !error ? (
        <span className={cx(fieldHintClassLayers.l1, fieldHintClassLayers.l2, fieldHintClassLayers.l3)}>{hint}</span>
      ) : null}
      {error ? (
        <span className={cx(fieldErrorClassLayers.l1, fieldErrorClassLayers.l2, fieldErrorClassLayers.l3)}>{error}</span>
      ) : null}
    </label>
  );
}

export function TextInput({ className = '', ...props }) {
  return (
    <input
      className={cx(
        inputClassLayers.l1,
        inputClassLayers.l2,
        inputClassLayers.l3,
        className
      )}
      {...props}
    />
  );
}
