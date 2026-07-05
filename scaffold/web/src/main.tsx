import React, { useCallback, useEffect, useState } from 'react';
import { createRoot } from 'react-dom/client';
import { AuthView, DashboardView, LoadingView } from './component/l3';
import type { AuthMode, AuthPayload, AuthResponse, MeResponse, ResolvedTheme, ThemeMode, User } from './lib/types';
import './tailwind.css';

const TEMPLATE_APP_NAME = '__' + 'PROJECT_NAME' + '__';
const PROJECT_APP_NAME = '__PROJECT_NAME__';
const APP_NAME = PROJECT_APP_NAME === TEMPLATE_APP_NAME ? 'Lorem Ipsum' : PROJECT_APP_NAME;
const THEME_STORAGE_KEY = 'carbide.theme';
const THEME_MODES = ['light', 'dark', 'system'] as const satisfies readonly ThemeMode[];

function isThemeMode(value: unknown): value is ThemeMode {
  return typeof value === 'string' && (THEME_MODES as readonly string[]).includes(value);
}

function systemTheme(): ResolvedTheme {
  if (!window.matchMedia) {
    return 'light';
  }
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

function storedThemeMode(): ThemeMode {
  try {
    const value = window.localStorage.getItem(THEME_STORAGE_KEY);
    return isThemeMode(value) ? value : 'system';
  } catch {
    return 'system';
  }
}

function applyThemeMode(mode: ThemeMode, resolved: ResolvedTheme): void {
  document.documentElement.dataset.theme = resolved;
  document.documentElement.dataset.themeMode = mode;
  document.documentElement.style.colorScheme = resolved;
}

function useThemeMode(): { mode: ThemeMode; resolved: ResolvedTheme; setMode: (nextMode: ThemeMode) => void } {
  const [mode, setModeState] = useState(storedThemeMode);
  const [system, setSystem] = useState(systemTheme);
  const resolved = mode === 'system' ? system : mode;

  useEffect(() => {
    if (!window.matchMedia) {
      return undefined;
    }
    const query = window.matchMedia('(prefers-color-scheme: dark)');
    const onChange = () => setSystem(query.matches ? 'dark' : 'light');
    query.addEventListener('change', onChange);
    return () => query.removeEventListener('change', onChange);
  }, []);

  useEffect(() => {
    applyThemeMode(mode, resolved);
    try {
      window.localStorage.setItem(THEME_STORAGE_KEY, mode);
    } catch {
      // Ignore storage failures; theme state still works for the current tab.
    }
  }, [mode, resolved]);

  const setMode = (nextMode: ThemeMode) => {
    if (isThemeMode(nextMode)) {
      setModeState(nextMode);
    }
  };

  return { mode, resolved, setMode };
}

async function api<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers = new Headers(options.headers);
  if (options.body && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/x-www-form-urlencoded');
  }

  const response = await fetch(path, {
    ...options,
    credentials: 'include',
    headers
  });
  const data = await response.json() as T & { error?: string };
  if (!response.ok) {
    throw new Error(data.error || 'Request failed.');
  }
  return data;
}

function encodeForm(values: Record<string, string>): string {
  const params = new URLSearchParams();
  Object.entries(values).forEach(([key, value]) => params.set(key, value));
  return params.toString();
}

function useRoute(): [string, (next: string) => void] {
  const [route, setRouteState] = useState(window.location.pathname);

  useEffect(() => {
    const onPop = () => setRouteState(window.location.pathname);
    window.addEventListener('popstate', onPop);
    return () => window.removeEventListener('popstate', onPop);
  }, []);

  const setRoute = useCallback((next: string) => {
    if (window.location.pathname !== next) {
      window.history.pushState({}, '', next);
    }
    setRouteState(next);
  }, []);

  return [route, setRoute];
}

function App() {
  const [route, setRoute] = useRoute();
  const theme = useThemeMode();
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');

  const mode: AuthMode = route === '/login' ? 'login' : 'register';

  useEffect(() => {
    api<MeResponse>('/api/me')
      .then((data) => {
        if (data.authenticated) {
          setUser(data.user);
          if (window.location.pathname === '/' || window.location.pathname === '/login') {
            setRoute('/dashboard');
          }
        } else if (window.location.pathname === '/dashboard') {
          setRoute('/login');
        }
      })
      .finally(() => setLoading(false));
  }, [setRoute]);

  const submitAuth = async ({ email, password }: AuthPayload) => {
    setBusy(true);
    setError('');
    try {
      const data = await api<AuthResponse>(`/api/${mode}`, {
        method: 'POST',
        body: encodeForm({ email, password })
      });
      setUser(data.user);
      setRoute('/dashboard');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Request failed.');
    } finally {
      setBusy(false);
    }
  };

  const logout = async () => {
    setBusy(true);
    try {
      await api<Record<string, never>>('/api/logout', { method: 'POST' });
      setUser(null);
      setRoute('/login');
    } finally {
      setBusy(false);
    }
  };

  if (loading) {
    return <LoadingView />;
  }

  if (route === '/dashboard' && user) {
    return (
      <DashboardView
        appName={APP_NAME}
        user={user}
        onLogout={logout}
        busy={busy}
        onThemeMode={theme.setMode}
        resolvedTheme={theme.resolved}
        themeMode={theme.mode}
      />
    );
  }

  return (
    <AuthView
      appName={APP_NAME}
      mode={mode}
      onSubmit={submitAuth}
      busy={busy}
      error={error}
      onThemeMode={theme.setMode}
      resolvedTheme={theme.resolved}
      themeMode={theme.mode}
      onMode={(nextMode) => {
        setError('');
        setRoute(nextMode === 'register' ? '/register' : '/login');
      }}
    />
  );
}

const rootElement = document.getElementById('root');
if (!rootElement) {
  throw new Error('missing #root element');
}

createRoot(rootElement).render(<App />);
