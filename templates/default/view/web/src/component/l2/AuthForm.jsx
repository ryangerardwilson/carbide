import { useEffect, useState } from 'react';
import { Button, Field, Muted, TextInput } from '../l1/index.js';

export function AuthForm({ busy, error = '', mode, onMode, onSubmit }) {
  const isRegister = mode === 'register';
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  useEffect(() => {
    setPassword('');
  }, [mode]);

  return (
    <form
      className="grid content-center gap-5 border-l border-emerald-950/10 bg-[#fbfdfb] px-7 py-10 sm:px-10 lg:min-h-svh lg:px-14"
      onSubmit={(event) => {
        event.preventDefault();
        onSubmit({ email, password });
      }}
    >
      {error ? <p className="m-0 rounded-md bg-rose-50 px-3 py-2 text-rose-800">{error}</p> : null}

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
          className="inline bg-transparent p-0 font-bold text-teal-700 underline-offset-4 hover:underline"
          type="button"
          onClick={() => onMode(isRegister ? 'login' : 'register')}
        >
          {isRegister ? 'Log in' : 'Create one'}
        </button>
      </Muted>
    </form>
  );
}
