import React from 'react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, waitFor, act } from '@testing-library/react';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import { AuthProvider, useAuth } from '../contexts/AuthContext';
import { apiClient } from '../api/client';
import { ProtectedRoute } from '../components/ProtectedRoute';

const ProtectedContent = () => <div>Protected Content</div>;
const LoginPage = () => <div>Login Page</div>;

describe('401 Response → Redirect Flow', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('full flow: 401 → onUnauthorized → logout → redirect (3.1.1-3.1.4)', async () => {
    const futureExp = Math.floor((Date.now() + 3600000) / 1000);
    const base64urlEncode = (str: string) =>
      btoa(str)
        .replace(/=/g, '')
        .replace(/\+/g, '-')
        .replace(/\//g, '_');
    const header = base64urlEncode(JSON.stringify({ alg: 'HS256', typ: 'JWT' }));
    const payload = base64urlEncode(
      JSON.stringify({ exp: futureExp, role: 'admin', id: 'usr_123' })
    );
    const validToken = `${header}.${payload}.mock-signature`;

    localStorage.setItem('auth_token', validToken);
    localStorage.setItem(
      'auth_user',
      JSON.stringify({ id: '1', name: 'Test', email: 'test@test.com', role: 'admin' })
    );

    const { container } = render(
      <MemoryRouter initialEntries={['/admin/dashboard']}>
        <AuthProvider>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route
              path="/admin/dashboard"
              element={
                <ProtectedRoute role="admin">
                  <ProtectedContent />
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(container.textContent).toContain('Protected Content');
    });

    const callback = apiClient.getOnUnauthorized();
    expect(callback).toBeDefined();

    if (callback) {
      act(() => {
        callback();
      });
    }

    await waitFor(() => {
      expect(container.textContent).toContain('Login Page');
    });
  });
});
