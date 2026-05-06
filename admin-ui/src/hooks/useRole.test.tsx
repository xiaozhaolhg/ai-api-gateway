import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { AuthProvider } from '../contexts/AuthContext';
import { useRole } from './useRole';
import { createMockToken } from '../test/utils';

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

  it('returns admin role when user has admin role', async () => {
    localStorage.setItem('auth_token', createMockToken({ role: 'admin' }));
    localStorage.setItem(
      'auth_user',
      JSON.stringify({ id: '1', name: 'Admin', email: 'admin@test.com', role: 'admin' })
    );

    const { result } = renderHook(() => useRole(), {
      wrapper: AuthProvider,
    });

    await waitFor(() => {
      expect(result.current).toBe('admin');
    });
  });

  it('returns user role when user has user role', async () => {
    localStorage.setItem('auth_token', createMockToken({ role: 'user' }));
    localStorage.setItem(
      'auth_user',
      JSON.stringify({ id: '2', name: 'User', email: 'user@test.com', role: 'user' })
    );

    const { result } = renderHook(() => useRole(), {
      wrapper: AuthProvider,
    });

    await waitFor(() => {
      expect(result.current).toBe('user');
    });
  });

  it('returns viewer role when user has viewer role', async () => {
    localStorage.setItem('auth_token', createMockToken({ role: 'viewer' }));
    localStorage.setItem(
      'auth_user',
      JSON.stringify({ id: '3', name: 'Viewer', email: 'viewer@test.com', role: 'viewer' })
    );

    const { result } = renderHook(() => useRole(), {
      wrapper: AuthProvider,
    });

    await waitFor(() => {
      expect(result.current).toBe('viewer');
    });
  });
});
