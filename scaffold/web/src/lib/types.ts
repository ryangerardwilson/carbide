export type AuthMode = 'login' | 'register';
export type ResolvedTheme = 'light' | 'dark';
export type ThemeMode = ResolvedTheme | 'system';

export interface AuthPayload {
  email: string;
  password: string;
}

export interface User {
  email: string;
}

export interface AuthResponse {
  user: User;
}

export type MeResponse =
  | {
      authenticated: true;
      user: User;
    }
  | {
      authenticated: false;
      user?: null;
    };

export interface NavItem {
  description?: string;
  eyebrow?: string;
  label: string;
  value: string;
}
