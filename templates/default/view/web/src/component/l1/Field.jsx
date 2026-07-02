import { cx } from '../utils.js';

export function Field({ children, error = '', hint = '', label, className = '' }) {
  return (
    <label className={cx('grid gap-2 font-bold text-[#16211b]', className)}>
      <span>{label}</span>
      {children}
      {hint && !error ? <span className="text-sm font-normal text-[#66786e]">{hint}</span> : null}
      {error ? <span className="text-sm font-bold text-rose-800">{error}</span> : null}
    </label>
  );
}

export function TextInput({ className = '', ...props }) {
  return (
    <input
      className={cx(
        'min-h-12 w-full rounded-md border border-emerald-900/20 bg-white px-3 py-2 text-[#16211b] outline-none transition focus:border-teal-700 focus:ring-4 focus:ring-teal-700/15',
        className
      )}
      {...props}
    />
  );
}
