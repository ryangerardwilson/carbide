import { AuthForm, LandingPageLayout } from '../l2';
import type { AuthMode, AuthPayload, ResolvedTheme, ThemeMode } from '../../lib/types';

interface AuthViewProps {
  appName: string;
  busy: boolean;
  error: string;
  mode: AuthMode;
  onMode: (mode: AuthMode) => void;
  onSubmit: (payload: AuthPayload) => void | Promise<void>;
  onThemeMode: (mode: ThemeMode) => void;
  resolvedTheme: ResolvedTheme;
  themeMode: ThemeMode;
}

export function AuthView({ appName, busy, error, mode, onMode, onSubmit, onThemeMode, resolvedTheme, themeMode }: AuthViewProps) {
  return (
    <LandingPageLayout appName={appName} onThemeMode={onThemeMode} resolvedTheme={resolvedTheme} themeMode={themeMode}>
      <AuthForm busy={busy} error={error} mode={mode} onMode={onMode} onSubmit={onSubmit} />
    </LandingPageLayout>
  );
}
