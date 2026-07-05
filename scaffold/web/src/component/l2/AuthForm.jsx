import { useEffect, useState } from 'react';
import { Button } from '../l1/Button.jsx';
import { Field, TextInput } from '../l1/Field.jsx';
import { Muted } from '../l1/Text.jsx';
import { ui } from '../l1/tokens.js';
import { cx } from '../../lib/cx.js';

const formClassLayers = {
  l1: 'grid content-center',
  l2: 'gap-3 border-l px-4 py-5 sm:px-5 lg:min-h-svh lg:px-8',
  l3: cx(ui.border, ui.surfaceSoft)
};

const formStackClassLayers = {
  l1: 'grid',
  l2: 'w-full max-w-sm justify-self-center gap-3',
  l3: ''
};

const errorClassLayers = {
  l1: '',
  l2: 'm-0 rounded-md px-2.5 py-1.5 text-sm/6',
  l3: cx(ui.errorSurface, ui.errorText)
};

const modeButtonClassLayers = {
  l1: 'inline',
  l2: 'rounded-sm p-0 font-semibold underline-offset-4',
  l3: cx('bg-transparent', ui.accent, ui.focus, 'hover:underline')
};

export function AuthForm({ busy, error = '', mode, onMode, onSubmit }) {
  const isRegister = mode === 'register';
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  useEffect(() => {
    setPassword('');
  }, [mode]);

  return (
    <form
      className={cx(formClassLayers.l1, formClassLayers.l2, formClassLayers.l3)}
      onSubmit={(event) => {
        event.preventDefault();
        onSubmit({ email, password });
      }}
    >
      <div className={cx(formStackClassLayers.l1, formStackClassLayers.l2, formStackClassLayers.l3)}>
        {error ? <p className={cx(errorClassLayers.l1, errorClassLayers.l2, errorClassLayers.l3)}>{error}</p> : null}

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
            className={cx(modeButtonClassLayers.l1, modeButtonClassLayers.l2, modeButtonClassLayers.l3)}
            type="button"
            onClick={() => onMode(isRegister ? 'login' : 'register')}
          >
            {isRegister ? 'Log in' : 'Create one'}
          </button>
        </Muted>
      </div>
    </form>
  );
}
