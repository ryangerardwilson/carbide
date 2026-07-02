import { cx } from '../utils.js';

export function Eyebrow({ children, className = '' }) {
  return (
    <p className={cx('m-0 text-xs font-extrabold uppercase tracking-normal text-teal-700', className)}>
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

  return <Tag className={cx('m-0 text-[#16211b]', sizes[level] || sizes[3], className)}>{children}</Tag>;
}

export function Muted({ children, className = '', as: Tag = 'p' }) {
  return <Tag className={cx('m-0 text-[#5d6f64]', className)}>{children}</Tag>;
}

export function CodeText({ children }) {
  return <code className="rounded bg-emerald-50 px-1.5 py-0.5 text-[#21463f]">{children}</code>;
}
