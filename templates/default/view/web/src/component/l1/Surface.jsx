import { cx } from '../utils.js';

export function Panel({ children, className = '', as: Tag = 'section', ...props }) {
  return (
    <Tag
      className={cx('rounded-lg border border-emerald-950/10 bg-white p-5 shadow-sm shadow-emerald-950/5', className)}
      {...props}
    >
      {children}
    </Tag>
  );
}

export function Divider({ className = '' }) {
  return <div className={cx('h-px w-full bg-emerald-950/10', className)} aria-hidden="true" />;
}

export function Badge({ children, tone = 'neutral', className = '' }) {
  const tones = {
    neutral: 'bg-stone-100 text-stone-700',
    good: 'bg-emerald-50 text-emerald-800',
    warn: 'bg-amber-50 text-amber-800',
    danger: 'bg-rose-50 text-rose-800'
  };

  return (
    <span className={cx('inline-flex min-h-7 items-center rounded-md px-2.5 text-sm font-bold', tones[tone], className)}>
      {children}
    </span>
  );
}

export function Metric({ label, value, detail = '' }) {
  return (
    <div className="min-w-0">
      <span className="mb-1 block text-sm text-[#6b7e72]">{label}</span>
      <strong className="block truncate text-[#16211b]">{value}</strong>
      {detail ? <span className="mt-1 block text-sm text-[#66786e]">{detail}</span> : null}
    </div>
  );
}
