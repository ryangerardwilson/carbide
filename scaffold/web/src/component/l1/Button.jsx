import { cx } from '../../lib/cx.js';
import { ui } from './tokens.js';

const variants = {
  primary: ui.action,
  secondary: `border ${ui.secondaryAction}`,
  ghost: ui.ghostAction,
  danger: ui.dangerAction
};

const sizes = {
  sm: 'min-h-7 px-2 text-xs',
  md: 'min-h-8 px-3 text-xs',
  lg: 'min-h-9 px-3.5 text-sm'
};

const buttonClassLayers = {
  l1: 'inline-flex items-center justify-center',
  l2: 'gap-1.5 rounded-md font-semibold outline-none',
  l3: 'transition disabled:cursor-wait disabled:opacity-65'
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
        buttonClassLayers.l1,
        buttonClassLayers.l2,
        sizes[size] || sizes.md,
        variants[variant] || variants.primary,
        buttonClassLayers.l3,
        className
      )}
      type={type}
      {...props}
    >
      {children}
    </button>
  );
}
