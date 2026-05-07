import React from 'react';
import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ProtectedRoute } from '../ProtectedRoute';
import { AuthProvider } from '../../contexts/AuthContext';
import { createMockToken, clearMockAuth } from '../../test/utils';

const renderWithRouter = (component: React.ReactElement) => {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });
  return render(
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <AuthProvider>
          {component}
        </AuthProvider>
      </BrowserRouter>
    </QueryClientProvider>
  );
};

describe('ProtectedRoute', () => {
  beforeEach(() => {
    clearMockAuth();
  });

  afterEach(() => {
    clearMockAuth();
  });

  it('should not render children when not authenticated', () => {
    renderWithRouter(<ProtectedRoute><div data-testid="protected">Protected</div></ProtectedRoute>);
    const protectedElement = screen.queryByTestId('protected');
    expect(protectedElement).not.toBeInTheDocument();
  });

  it('should render children when authenticated', async () => {
    localStorage.setItem('auth_token', createMockToken({ role: 'admin' }));
    localStorage.setItem('auth_user', JSON.stringify({
      id: '1',
      name: 'Test User',
      email: 'test@example.com',
      role: 'admin',
    }));

    renderWithRouter(
      <ProtectedRoute>
        <div data-testid="protected">Protected Content</div>
      </ProtectedRoute>
    );

    await waitFor(() => {
      expect(screen.getByTestId('protected')).toBeInTheDocument();
    });
    expect(screen.getByText('Protected Content')).toBeInTheDocument();
  });

  it('should respect requiredRole for admin', async () => {
    localStorage.setItem('auth_token', createMockToken({ role: 'admin' }));
    localStorage.setItem('auth_user', JSON.stringify({
      id: '1',
      name: 'Test User',
      email: 'test@example.com',
      role: 'admin',
    }));

    renderWithRouter(
      <ProtectedRoute requiredRole="admin">
        <div data-testid="protected">Admin Only</div>
      </ProtectedRoute>
    );

    await waitFor(() => {
      expect(screen.getByTestId('protected')).toBeInTheDocument();
    });
  });
});
