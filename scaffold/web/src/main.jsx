import React, { useEffect, useMemo, useState } from 'react';
import { createRoot } from 'react-dom/client';
import { AuthView, DashboardView, LoadingView } from './component/l3/index.js';
import './tailwind.css';

const APP_NAME = '__PROJECT_NAME__';

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
    return <DashboardView appName={APP_NAME} user={user} onLogout={logout} busy={busy} />;
  }

  return (
    <AuthView
      appName={APP_NAME}
      mode={mode}
      onSubmit={submitAuth}
      busy={busy}
      error={error}
      onMode={(nextMode) => {
        setError('');
        setRoute(nextMode === 'register' ? '/register' : '/login');
      }}
    />
  );
}

createRoot(document.getElementById('root')).render(<App />);
