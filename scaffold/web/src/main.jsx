import React, { useEffect, useMemo, useState } from 'react';
import { createRoot } from 'react-dom/client';
import { AuthView, DashboardView, LoadingView } from './component/l3/index.js';
import './tailwind.css';

const TEMPLATE_APP_NAME = '__' + 'PROJECT_NAME' + '__';
const PROJECT_APP_NAME = '__PROJECT_NAME__';
const APP_NAME = PROJECT_APP_NAME === TEMPLATE_APP_NAME ? 'Lorem Ipsum' : PROJECT_APP_NAME;
const THEME_STORAGE_KEY = 'carbide.theme';
const THEME_MODES = ['light', 'dark', 'system'];

function systemTheme() {
  if (!window.matchMedia) {
    return 'light';
  }
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

function storedThemeMode() {
  try {
    const value = window.localStorage.getItem(THEME_STORAGE_KEY);
    return THEME_MODES.includes(value) ? value : 'system';
  } catch {
    return 'system';
  }
}

function applyThemeMode(mode, resolved) {
  document.documentElement.dataset.theme = resolved;
  document.documentElement.dataset.themeMode = mode;
  document.documentElement.style.colorScheme = resolved;
}

function useThemeMode() {
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

  const setMode = (nextMode) => {
    if (THEME_MODES.includes(nextMode)) {
      setModeState(nextMode);
    }
  };

  return { mode, resolved, setMode };
}

async function api(path, options = {}) {
  const response = await fetch(path, {
    credentials: 'include',
    ...options,
    headers: {
      ...(options.body ? { 'Content-Type': 'application/x-www-form-urlencoded' } : {}),
      ...(options.headers || {})
    }
  });
  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.error || 'Request failed.');
  }
  return data;
}

function encodeForm(values) {
  const params = new URLSearchParams();
  Object.entries(values).forEach(([key, value]) => params.set(key, value));
  return params.toString();
}

function useRoute() {
  const [route, setRouteState] = useState(window.location.pathname);

  useEffect(() => {
    const onPop = () => setRouteState(window.location.pathname);
    window.addEventListener('popstate', onPop);
    return () => window.removeEventListener('popstate', onPop);
  }, []);

  const setRoute = (next) => {
    if (window.location.pathname !== next) {
      window.history.pushState({}, '', next);
    }
    setRouteState(next);
  };

  return [route, setRoute];
}

function App() {
  const [route, setRoute] = useRoute();
  const theme = useThemeMode();
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');

  const mode = useMemo(() => (route === '/login' ? 'login' : 'register'), [route]);

  useEffect(() => {
    api('/api/me')
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
  }, []);

  const submitAuth = async ({ email, password }) => {
    setBusy(true);
    setError('');
    try {
      const data = await api(`/api/${mode}`, {
        method: 'POST',
        body: encodeForm({ email, password })
      });
      setUser(data.user);
      setRoute('/dashboard');
    } catch (err) {
      setError(err.message);
    } finally {
      setBusy(false);
    }
  };

  const logout = async () => {
    setBusy(true);
    try {
      await api('/api/logout', { method: 'POST' });
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

createRoot(document.getElementById('root')).render(<App />);
