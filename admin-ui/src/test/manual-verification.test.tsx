import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { AuthProvider, useAuth } from '../contexts/AuthContext';
import { ProtectedRoute } from '../components/ProtectedRoute';
import { createMockToken } from './utils';

describe('JWT Expiration Manual Verification', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  afterEach(() => {
    localStorage.clear();
  });

  it('should initialize AuthContext with stored credentials', async () => {
    localStorage.setItem('auth_token', createMockToken({ role: 'admin' }));
    localStorage.setItem('auth_user', JSON.stringify({
      id: '1',
      name: 'Test User',
      email: 'test@example.com',
      role: 'admin',
    }));

    const TestComponent = () => {
      const { user, isAuthenticated } = useAuth();

      return (
        <div>
          <p data-testid="auth-status">Authenticated: {isAuthenticated ? 'true' : 'false'}</p>
          <p data-testid="user-email">User: {user?.email || 'null'}</p>
          <p data-testid="user-role">Role: {user?.role || 'null'}</p>
        </div>
      );
    };

    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    );

    await waitFor(() => {
      expect(screen.getByTestId('auth-status')).toHaveTextContent('Authenticated: true');
    });
    expect(screen.getByTestId('user-email')).toHaveTextContent('User: test@example.com');
    expect(screen.getByTestId('user-role')).toHaveTextContent('Role: admin');
  });

  it('should handle ProtectedRoute with admin role', async () => {
    localStorage.setItem('auth_token', createMockToken({ role: 'admin' }));
    localStorage.setItem('auth_user', JSON.stringify({
      id: '1',
      name: 'Test User',
      email: 'test@example.com',
      role: 'admin',
    }));

    render(
      <MemoryRouter>
        <AuthProvider>
          <ProtectedRoute requiredRole="admin">
            <div data-testid="protected-content">Admin Content</div>
          </ProtectedRoute>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByTestId('protected-content')).toBeInTheDocument();
    });
    expect(screen.getByTestId('protected-content')).toHaveTextContent('Admin Content');
  });

  it('should handle ProtectedRoute with user role for admin user', async () => {
    localStorage.setItem('auth_token', createMockToken({ role: 'admin' }));
    localStorage.setItem('auth_user', JSON.stringify({
      id: '1',
      name: 'Test User',
      email: 'test@example.com',
      role: 'admin',
    }));

    render(
      <MemoryRouter>
        <AuthProvider>
          <ProtectedRoute requiredRole="user">
            <div data-testid="protected-content">User Content</div>
          </ProtectedRoute>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByTestId('protected-content')).toBeInTheDocument();
    });
    expect(screen.getByTestId('protected-content')).toHaveTextContent('User Content');
  });
});
