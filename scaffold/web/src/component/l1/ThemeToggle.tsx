import type { SVGProps } from 'react';
import { cx } from '../../lib/cx';
import type { ResolvedTheme, ThemeMode } from '../../lib/types';
import { ui } from './tokens';

const themeToggleClassLayers = {
  button: {
    l1: 'inline-flex items-center justify-center',
    l2: 'size-8 rounded-full border outline-none',
    l3: cx(ui.secondaryAction, ui.focus, 'transition hover:border-carbide-border-strong')
  },
  icon: {
    l1: 'shrink-0',
    l2: 'size-4',
    l3: ''
  }
};

function SunIcon({ className = '' }: SVGProps<SVGSVGElement>) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" aria-hidden="true">
      <circle cx="12" cy="12" r="3.5" stroke="currentColor" strokeWidth="2" />
      <path
        d="M12 2.75v2.5M12 18.75v2.5M4.42 4.42l1.77 1.77M17.81 17.81l1.77 1.77M2.75 12h2.5M18.75 12h2.5M4.42 19.58l1.77-1.77M17.81 6.19l1.77-1.77"
        stroke="currentColor"
        strokeLinecap="round"
        strokeWidth="2"
      />
    </svg>
  );
}

function MoonIcon({ className = '' }: SVGProps<SVGSVGElement>) {
  return (
    <svg className={className} fill="none" viewBox="0 0 24 24" aria-hidden="true">
      <path
        d="M20.25 14.42A7.88 7.88 0 0 1 9.58 3.75 8.5 8.5 0 1 0 20.25 14.42Z"
        stroke="currentColor"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="2"
      />
    </svg>
  );
}

interface ThemeToggleProps {
  className?: string;
  mode?: ThemeMode;
  onMode?: (mode: ThemeMode) => void;
  resolved?: ResolvedTheme;
}

export function ThemeToggle({ className = '', mode = 'system', onMode, resolved = 'light' }: ThemeToggleProps) {
  const isDark = resolved === 'dark';
  const nextMode = isDark ? 'light' : 'dark';
  const label = isDark ? 'Switch to light theme' : 'Switch to dark theme';
  const Icon = isDark ? MoonIcon : SunIcon;

  return (
    <button
      aria-label={label}
      className={cx(themeToggleClassLayers.button.l1, themeToggleClassLayers.button.l2, themeToggleClassLayers.button.l3, className)}
      data-resolved-theme={resolved}
      data-theme-mode={mode}
      type="button"
      onClick={() => onMode?.(nextMode)}
    >
      <Icon className={cx(themeToggleClassLayers.icon.l1, themeToggleClassLayers.icon.l2, themeToggleClassLayers.icon.l3)} />
    </button>
  );
}
