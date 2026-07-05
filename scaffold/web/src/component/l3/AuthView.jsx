import { AuthForm, LandingPageLayout } from '../l2/index.js';

export function AuthView({ appName, busy, error, mode, onMode, onSubmit, onThemeMode, resolvedTheme, themeMode }) {
  return (
    <LandingPageLayout appName={appName} onThemeMode={onThemeMode} resolvedTheme={resolvedTheme} themeMode={themeMode}>
      <AuthForm busy={busy} error={error} mode={mode} onMode={onMode} onSubmit={onSubmit} />
    </LandingPageLayout>
  );
}
