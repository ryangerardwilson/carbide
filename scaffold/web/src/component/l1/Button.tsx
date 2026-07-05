import type { ButtonHTMLAttributes, ReactNode } from 'react';
import { cx } from '../../lib/cx';
import { ui } from './tokens';

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

type ButtonSize = keyof typeof sizes;
type ButtonVariant = keyof typeof variants;

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  children: ReactNode;
  className?: string;
  size?: ButtonSize;
  variant?: ButtonVariant;
}

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
}: ButtonProps) {
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
