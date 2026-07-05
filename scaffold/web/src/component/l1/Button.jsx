import { cx } from '../../lib/cx.js';
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

const buttonClassLayers = {
  l1: 'inline-flex items-center justify-center',
  l2: 'gap-2 rounded-md font-bold outline-none',
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
