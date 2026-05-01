import React, { ReactNode } from 'react';
import { AuthProvider } from '../contexts/AuthContext';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

export const createMockUser = (overrides: Partial<User> = {}): User => ({
  id: 'usr_test123',
  name: 'Test User',
  email: 'test@example.com',
  role: 'admin' as const,
  status: 'active',
  created_at: new Date().toISOString(),
  ...overrides,
});

export interface User {
  id: string;
  name: string;
  email: string;
  role: 'admin' | 'user' | 'viewer';
  status: string;
  created_at: string;
}

export const MockAuthProvider: React.FC<{
  children: ReactNode;
  role?: 'admin' | 'user' | 'viewer';
  isAuthenticated?: boolean;
}> = ({ children, role = 'admin', isAuthenticated = true }) => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

  const mockUser = createMockUser({ role });

  const mockAuthValue = {
    user: isAuthenticated ? mockUser : null,
    isAuthenticated,
    login: async () => ({ token: 'mock-token', user: mockUser }),
    logout: async () => {},
    getCurrentUser: async () => mockUser,
  };

  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider value={mockAuthValue}>
        {children}
      </AuthProvider>
    </QueryClientProvider>
  );
};

export const renderWithAuth = (
  ui: React.ReactElement,
  { role = 'admin', isAuthenticated = true } = {}
) => {
  return render(ui, {
    wrapper: ({ children }) => (
      <MockAuthProvider role={role} isAuthenticated={isAuthenticated}>
        {children}
      </MockAuthProvider>
    ),
  });
};

export const createMockToken = (role: 'admin' | 'user' | 'viewer' = 'admin') => {
  return `mock-jwt-token-${role}-${Date.now()}`;
};

export const setMockAuthToken = (token: string) => {
  localStorage.setItem('auth_token', token);
};

export const clearMockAuth = () => {
  localStorage.removeItem('auth_token');
};
