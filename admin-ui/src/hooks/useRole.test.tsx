import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook } from '@testing-library/react';
import { AuthProvider } from '../contexts/AuthContext';
import { useRole } from './useRole';

describe('useRole', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  it('returns viewer as default when no user is authenticated', () => {
    const { result } = renderHook(() => useRole(), {
      wrapper: AuthProvider,
    });

    expect(result.current).toBe('viewer');
  });

  it('returns admin role when user has admin role', () => {
    localStorage.setItem('auth_token', 'test-token');
    localStorage.setItem(
      'auth_user',
      JSON.stringify({ id: '1', name: 'Admin', email: 'admin@test.com', role: 'admin' })
    );

    const { result } = renderHook(() => useRole(), {
      wrapper: AuthProvider,
    });

    expect(result.current).toBe('admin');
  });

  it('returns user role when user has user role', () => {
    localStorage.setItem('auth_token', 'test-token');
    localStorage.setItem(
      'auth_user',
      JSON.stringify({ id: '2', name: 'User', email: 'user@test.com', role: 'user' })
    );

    const { result } = renderHook(() => useRole(), {
      wrapper: AuthProvider,
    });

    expect(result.current).toBe('user');
  });

  it('returns viewer role when user has viewer role', () => {
    localStorage.setItem('auth_token', 'test-token');
    localStorage.setItem(
      'auth_user',
      JSON.stringify({ id: '3', name: 'Viewer', email: 'viewer@test.com', role: 'viewer' })
    );

    const { result } = renderHook(() => useRole(), {
      wrapper: AuthProvider,
    });

    expect(result.current).toBe('viewer');
  });
});
