import { cx } from '../utils.js';

const variants = {
  primary: 'bg-teal-700 text-white hover:bg-teal-800 focus-visible:ring-teal-700/25',
  secondary:
    'border border-emerald-950/15 bg-white text-[#16211b] hover:border-emerald-950/30 focus-visible:ring-emerald-950/15',
  ghost: 'bg-transparent text-teal-700 hover:bg-teal-50 focus-visible:ring-teal-700/15',
  danger: 'bg-rose-700 text-white hover:bg-rose-800 focus-visible:ring-rose-700/20'
};

const sizes = {
  sm: 'min-h-9 px-3 text-sm',
  md: 'min-h-11 px-5',
  lg: 'min-h-12 px-6 text-lg'
};

export function Button({
  children,
  className = '',
  size = 'md',
  type = 'button',
  variant = 'primary',
  ...props
}) {
  return (
    <button
      className={cx(
        'inline-flex items-center justify-center gap-2 rounded-md font-bold outline-none transition focus-visible:ring-4 disabled:cursor-wait disabled:opacity-65',
        variants[variant] || variants.primary,
        sizes[size] || sizes.md,
        className
      )}
      type={type}
      {...props}
    >
      {children}
    </button>
  );
}
