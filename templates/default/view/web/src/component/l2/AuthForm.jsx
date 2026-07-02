import { useEffect, useState } from 'react';
import { Button, Field, Muted, TextInput, ui } from '../l1/index.js';
import { cx } from '../utils.js';

export function AuthForm({ busy, error = '', mode, onMode, onSubmit }) {
  const isRegister = mode === 'register';
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  useEffect(() => {
    setPassword('');
  }, [mode]);

  return (
    <form
      className={cx('grid content-center gap-5 border-l px-7 py-10 sm:px-10 lg:min-h-svh lg:px-14', ui.border, ui.surfaceSoft)}
      onSubmit={(event) => {
        event.preventDefault();
        onSubmit({ email, password });
      }}
    >
      {error ? <p className={cx('m-0 rounded-md px-3 py-2', ui.errorSurface, ui.errorText)}>{error}</p> : null}

      <Field label="Email">
        <TextInput
          autoComplete="email"
          name="email"
          onChange={(event) => setEmail(event.target.value)}
          required
          type="email"
          value={email}
        />
      </Field>

      <Field label="Password">
        <TextInput
          autoComplete={isRegister ? 'new-password' : 'current-password'}
          name="password"
          onChange={(event) => setPassword(event.target.value)}
          required
          type="password"
          value={password}
        />
      </Field>

      <Button disabled={busy} type="submit">
        {busy ? 'Working...' : isRegister ? 'Create account' : 'Log in'}
      </Button>

      <Muted>
        {isRegister ? 'Already registered?' : 'Need an account?'}{' '}
        <button
          className={cx('inline bg-transparent p-0 font-bold underline-offset-4 hover:underline', ui.accent)}
          type="button"
          onClick={() => onMode(isRegister ? 'login' : 'register')}
        >
          {isRegister ? 'Log in' : 'Create one'}
        </button>
      </Muted>
    </form>
  );
}
