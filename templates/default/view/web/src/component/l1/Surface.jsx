import { cx } from '../utils.js';
import { ui } from './tokens.js';

export function Panel({ children, className = '', as: Tag = 'section', ...props }) {
  return (
    <Tag
      className={cx('rounded-lg border p-5', ui.border, ui.surface, ui.shadowSubtle, className)}
      {...props}
    >
      {children}
    </Tag>
  );
}

export function Divider({ className = '' }) {
  return <div className={cx('h-px w-full', ui.divider, className)} aria-hidden="true" />;
}

export function Badge({ children, tone = 'neutral', className = '' }) {
  const tones = {
    neutral: ui.neutralBadge,
    good: ui.goodBadge,
    warn: ui.warnBadge,
    danger: ui.dangerBadge
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
      <span className={cx('mb-1 block text-sm', ui.subtle)}>{label}</span>
      <strong className={cx('block truncate', ui.text)}>{value}</strong>
      {detail ? <span className={cx('mt-1 block text-sm', ui.subtle)}>{detail}</span> : null}
    </div>
  );
}
