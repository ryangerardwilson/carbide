import type { ElementType, ReactNode } from 'react';
import { cx } from '../../lib/cx';
import { ui } from './tokens';

const eyebrowClassLayers = {
  l1: '',
  l2: 'm-0 text-xs font-bold uppercase tracking-normal',
  l3: ui.accent
};

const headingClassLayers = {
  l1: '',
  l2: 'm-0',
  l3: ui.text
};

const mutedClassLayers = {
  l1: '',
  l2: 'm-0 text-sm/6',
  l3: ui.muted
};

const codeClassLayers = {
  l1: '',
  l2: 'rounded px-1 py-0.5 text-xs',
  l3: ui.code
};

interface TextProps {
  children: ReactNode;
  className?: string;
}

interface HeadingProps extends TextProps {
  level?: 1 | 2 | 3;
}

interface MutedProps extends TextProps {
  as?: ElementType;
}

export function Eyebrow({ children, className = '' }: TextProps) {
  return (
    <p className={cx(eyebrowClassLayers.l1, eyebrowClassLayers.l2, eyebrowClassLayers.l3, className)}>
      {children}
    </p>
  );
}

export function Heading({ children, className = '', level = 1 }: HeadingProps) {
  const Tag = `h${level}` as ElementType;
  const sizes = {
    1: 'text-2xl/8 sm:text-3xl/9',
    2: 'text-xl/7',
    3: 'text-base/6'
  };

  return <Tag className={cx(headingClassLayers.l1, headingClassLayers.l2, sizes[level] || sizes[3], headingClassLayers.l3, className)}>{children}</Tag>;
}

export function Muted({ children, className = '', as: Tag = 'p' }: MutedProps) {
  return <Tag className={cx(mutedClassLayers.l1, mutedClassLayers.l2, mutedClassLayers.l3, className)}>{children}</Tag>;
}

export function CodeText({ children }: { children: ReactNode }) {
  return <code className={cx(codeClassLayers.l1, codeClassLayers.l2, codeClassLayers.l3)}>{children}</code>;
}
