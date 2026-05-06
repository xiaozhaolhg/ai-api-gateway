import type { ReactNode } from 'react';
import { createContext, useContext } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render } from '@testing-library/react';

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

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  login: () => Promise<{ token: string; user: User }>;
  logout: () => Promise<void>;
  getCurrentUser: () => Promise<User>;
}

const MockAuthContext = createContext<AuthContextType | undefined>(undefined);

export const MockAuthProvider: React.FC<{
  children: ReactNode;
  role?: 'admin' | 'user' | 'viewer';
  isAuthenticated?: boolean;
}> = ({ children, role = 'admin', isAuthenticated = true }) => {
  const mockUser = createMockUser({ role });

  const mockAuthValue: AuthContextType = {
    user: isAuthenticated ? mockUser : null,
    isAuthenticated,
    login: async () => ({ token: 'mock-token', user: mockUser }),
    logout: async () => {},
    getCurrentUser: async () => mockUser,
  };

  return (
    <MockAuthContext.Provider value={mockAuthValue}>
      {children}
    </MockAuthContext.Provider>
  );
};

export const useMockAuth = () => {
  const context = useContext(MockAuthContext);
  if (!context) {
    throw new Error('useMockAuth must be used within MockAuthProvider');
  }
  return context;
};

export const renderWithAuth = (
  ui: React.ReactElement,
  { role = 'admin', isAuthenticated = true }: { role?: 'admin' | 'user' | 'viewer'; isAuthenticated?: boolean } = {}
) => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

  const Wrapper: React.FC<{ children: ReactNode }> = ({ children }) => (
    <QueryClientProvider client={queryClient}>
      <MockAuthProvider role={role} isAuthenticated={isAuthenticated}>
        {children}
      </MockAuthProvider>
    </QueryClientProvider>
  );

  return render(ui, { wrapper: Wrapper });
};

export const createMockToken = ({
  role = 'admin' as const,
  expiresIn = 86400000,
}: {
  role?: 'admin' | 'user' | 'viewer';
  expiresIn?: number;
} = {}) => {
  const base64urlEncode = (str: string) => {
    return btoa(str)
      .replace(/=/g, '')
      .replace(/\+/g, '-')
      .replace(/\//g, '_');
  };

  const header = base64urlEncode(JSON.stringify({ alg: 'HS256', typ: 'JWT' }));
  const exp = Math.floor((Date.now() + expiresIn) / 1000);
  const payload = base64urlEncode(JSON.stringify({ exp, role, id: 'usr_test123' }));
  const signature = 'mock-signature';
  return `${header}.${payload}.${signature}`;
};

export const setMockAuthToken = (token: string) => {
  localStorage.setItem('auth_token', token);
};

export const clearMockAuth = () => {
  localStorage.removeItem('auth_token');
};
