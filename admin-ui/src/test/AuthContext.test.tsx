import { describe, it, expect, vi, beforeEach, afterEach, beforeAll, afterAll } from 'vitest';
import { renderHook, act, waitFor } from '@testing-library/react';
import { AuthProvider, useAuth } from '../contexts/AuthContext';

const ORIGINAL_ENV = process.env.VITE_USE_MOCK;

beforeAll(() => {
  process.env.VITE_USE_MOCK = 'true';
});

afterAll(() => {
  process.env.VITE_USE_MOCK = ORIGINAL_ENV;
});

describe('AuthContext', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('provides initial null user and token', () => {
    const { result } = renderHook(() => useAuth(), {
      wrapper: AuthProvider,
    });

    expect(result.current.user).toBeNull();
    expect(result.current.token).toBeNull();
    expect(result.current.isAuthenticated).toBe(false);
  });

  it('initializes with stored credentials from localStorage', async () => {
    localStorage.setItem('auth_token', createFutureToken(3600000));
    localStorage.setItem(
      'auth_user',
      JSON.stringify({ id: '1', name: 'Stored', email: 'stored@test.com', role: 'admin' })
    );

    const { result, rerender } = renderHook(() => useAuth(), {
      wrapper: AuthProvider,
    });

    // Wait for AuthContext to initialize from localStorage
    await waitFor(() => {
      expect(result.current.isAuthenticated).toBe(true);
      expect(result.current.user?.email).toBe('stored@test.com');
    });
  });

  it('logout clears user and token', async () => {
    localStorage.setItem('auth_token', createFutureToken(3600000));
    localStorage.setItem(
      'auth_user',
      JSON.stringify({ id: '1', name: 'Test', email: 'test@test.com', role: 'admin' })
    );

    const { result } = renderHook(() => useAuth(), {
      wrapper: AuthProvider,
    });

    // Wait for AuthContext to initialize from localStorage
    await waitFor(() => {
      expect(result.current.isAuthenticated).toBe(true);
    });

    act(() => {
      result.current.logout();
    });

    expect(result.current.isAuthenticated).toBe(false);
    expect(result.current.user).toBeNull();
    expect(result.current.token).toBeNull();
  });

  it('throws error when useAuth used outside provider', () => {
    const consoleError = vi.spyOn(console, 'error').mockImplementation(() => {});

    expect(() => {
      renderHook(() => useAuth());
    }).toThrow('useAuth must be used within an AuthProvider');

    consoleError.mockRestore();
  });

  it('detects expired token on initialization', () => {
    const expiredToken = createExpiredToken();
    localStorage.setItem('auth_token', expiredToken);
    localStorage.setItem(
      'auth_user',
      JSON.stringify({ id: '1', name: 'Test', email: 'test@test.com', role: 'admin' })
    );

    const { result } = renderHook(() => useAuth(), {
      wrapper: AuthProvider,
    });

    expect(result.current.isAuthenticated).toBe(false);
  });

  it('does not logout with valid future token', () => {
    const futureToken = createFutureToken(3600000);
    localStorage.setItem('auth_token', futureToken);
    localStorage.setItem(
      'auth_user',
      JSON.stringify({ id: '1', name: 'Test', email: 'test@test.com', role: 'admin' })
    );

    const { result } = renderHook(() => useAuth(), {
      wrapper: AuthProvider,
    });

    expect(result.current.isAuthenticated).toBe(true);
  });

  it('handles token without exp claim gracefully', () => {
    const noExpToken = createTokenWithoutExp();
    localStorage.setItem('auth_token', noExpToken);
    localStorage.setItem(
      'auth_user',
      JSON.stringify({ id: '1', name: 'Test', email: 'test@test.com', role: 'admin' })
    );

    const { result } = renderHook(() => useAuth(), {
      wrapper: AuthProvider,
    });

    expect(result.current.isAuthenticated).toBe(true);
  });
});

const base64urlEncode = (str: string) => {
  return btoa(str)
    .replace(/=/g, '')
    .replace(/\+/g, '-')
    .replace(/\//g, '_');
};

const createExpiredToken = () => {
  const header = base64urlEncode(JSON.stringify({ alg: 'HS256', typ: 'JWT' }));
  const payload = base64urlEncode(
    JSON.stringify({ exp: Math.floor((Date.now() - 86400000) / 1000), role: 'admin', id: 'usr_123' })
  );
  return `${header}.${payload}.mock-signature`;
};

const createFutureToken = (ms: number) => {
  const header = base64urlEncode(JSON.stringify({ alg: 'HS256', typ: 'JWT' }));
  const payload = base64urlEncode(
    JSON.stringify({ exp: Math.floor((Date.now() + ms) / 1000), role: 'admin', id: 'usr_123' })
  );
  return `${header}.${payload}.mock-signature`;
};

const createTokenWithoutExp = () => {
  const header = base64urlEncode(JSON.stringify({ alg: 'HS256', typ: 'JWT' }));
  const payload = base64urlEncode(
    JSON.stringify({ role: 'admin', id: 'usr_123' })
  );
  return `${header}.${payload}.mock-signature`;
};