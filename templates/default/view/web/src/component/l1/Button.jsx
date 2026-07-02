import { cx } from '../utils.js';
import { ui } from './tokens.js';

const variants = {
  primary: ui.action,
  secondary: `border ${ui.secondaryAction}`,
  ghost: ui.ghostAction,
  danger: ui.dangerAction
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
        'inline-flex items-center justify-center gap-2 rounded-md font-bold outline-none transition disabled:cursor-wait disabled:opacity-65',
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
