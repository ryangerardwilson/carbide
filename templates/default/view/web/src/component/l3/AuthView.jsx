import { AuthForm, LandingPageLayout } from '../l2/index.js';

export function AuthView({ appName, busy, error, mode, onMode, onSubmit }) {
  return (
    <LandingPageLayout appName={appName} mode={mode}>
      <AuthForm busy={busy} error={error} mode={mode} onMode={onMode} onSubmit={onSubmit} />
    </LandingPageLayout>
  );
}
