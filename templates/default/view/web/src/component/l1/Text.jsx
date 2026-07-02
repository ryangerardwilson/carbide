import { cx } from '../utils.js';
import { ui } from './tokens.js';

export function Eyebrow({ children, className = '' }) {
  return (
    <p className={cx('m-0 text-xs font-extrabold uppercase tracking-normal', ui.accent, className)}>
      {children}
    </p>
  );
}

export function Heading({ children, className = '', level = 1 }) {
  const Tag = `h${level}`;
  const sizes = {
    1: 'text-[34px] leading-tight sm:text-5xl',
    2: 'text-3xl leading-tight',
    3: 'text-xl leading-snug'
  };

  return <Tag className={cx('m-0', ui.text, sizes[level] || sizes[3], className)}>{children}</Tag>;
}

export function Muted({ children, className = '', as: Tag = 'p' }) {
  return <Tag className={cx('m-0', ui.muted, className)}>{children}</Tag>;
}

export function CodeText({ children }) {
  return <code className={cx('rounded px-1.5 py-0.5', ui.code)}>{children}</code>;
}
