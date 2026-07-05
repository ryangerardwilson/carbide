import { cx } from '../../lib/cx.js';
import { ui } from './tokens.js';

const panelClassLayers = {
  l1: '',
  l2: 'rounded-lg border p-3',
  l3: cx(ui.border, ui.surface, ui.shadowSubtle)
};

const dividerClassLayers = {
  l1: 'w-full',
  l2: 'h-px',
  l3: ui.divider
};

const badgeClassLayers = {
  l1: 'inline-flex items-center',
  l2: 'min-h-6 rounded-md px-2 text-xs font-semibold'
};

const metricClassLayers = {
  root: {
    l1: 'min-w-0',
    l2: '',
    l3: ''
  },
  label: {
    l1: 'block',
    l2: 'mb-0.5 text-xs',
    l3: ui.subtle
  },
  value: {
    l1: 'block truncate',
    l2: 'text-sm/6 font-semibold',
    l3: ui.text
  },
  detail: {
    l1: 'block',
    l2: 'mt-1 text-xs',
    l3: ui.subtle
  }
};

export function Panel({ children, className = '', as: Tag = 'section', ...props }) {
  return (
    <Tag
      className={cx(panelClassLayers.l1, panelClassLayers.l2, panelClassLayers.l3, className)}
      {...props}
    >
      {children}
    </Tag>
  );
}

export function Divider({ className = '' }) {
  return <div className={cx(dividerClassLayers.l1, dividerClassLayers.l2, dividerClassLayers.l3, className)} aria-hidden="true" />;
}

export function Badge({ children, tone = 'neutral', className = '' }) {
  const tones = {
    neutral: ui.neutralBadge,
    good: ui.goodBadge,
    warn: ui.warnBadge,
    danger: ui.dangerBadge
  };

  return (
    <span className={cx(badgeClassLayers.l1, badgeClassLayers.l2, tones[tone], className)}>
      {children}
    </span>
  );
}

export function Metric({ label, value, detail = '' }) {
  return (
    <div className={cx(metricClassLayers.root.l1, metricClassLayers.root.l2, metricClassLayers.root.l3)}>
      <span className={cx(metricClassLayers.label.l1, metricClassLayers.label.l2, metricClassLayers.label.l3)}>
        {label}
      </span>
      <strong className={cx(metricClassLayers.value.l1, metricClassLayers.value.l2, metricClassLayers.value.l3)}>
        {value}
      </strong>
      {detail ? (
        <span className={cx(metricClassLayers.detail.l1, metricClassLayers.detail.l2, metricClassLayers.detail.l3)}>
          {detail}
        </span>
      ) : null}
    </div>
  );
}
