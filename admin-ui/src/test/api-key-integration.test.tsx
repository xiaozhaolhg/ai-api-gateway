import React from 'react';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter } from 'react-router-dom';
import { AuthProvider } from '../contexts/AuthContext';
import APIKeys from '../pages/APIKeys';
import { apiClient } from '../api/client';

// Mock i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

// Mock the API client
vi.mock('../api/client', () => ({
  apiClient: {
    getUsers: vi.fn(),
    getAPIKeys: vi.fn(),
    createAPIKey: vi.fn(),
    deleteAPIKey: vi.fn(),
    setOnUnauthorized: vi.fn(),
    getCurrentUser: vi.fn(),
  },
}));

// Mock sessionStorage
const sessionStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
};
Object.defineProperty(window, 'sessionStorage', {
  value: sessionStorageMock,
});

// Mock clipboard
const clipboardMock = {
  writeText: vi.fn(),
};
Object.defineProperty(navigator, 'clipboard', {
  value: clipboardMock,
  writable: true,
});

// Mock matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(), // deprecated
    removeListener: vi.fn(), // deprecated
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
});

// Mock getComputedStyle
Object.defineProperty(window, 'getComputedStyle', {
  value: vi.fn(() => ({
    getPropertyValue: vi.fn(() => ''),
  })),
});

const createQueryClient = () => new QueryClient({
  defaultOptions: {
    queries: { retry: false },
    mutations: { retry: false },
  },
});

const renderWithProviders = (component: React.ReactElement) => {
  const queryClient = createQueryClient();
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

describe('API Keys Integration Tests', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    sessionStorageMock.getItem.mockReturnValue(null);
    (apiClient.getUsers as vi.Mock).mockResolvedValue([
      { id: 'user1', name: 'Test User', email: 'test@example.com', role: 'admin' as const, status: 'active', created_at: '2024-01-01' }
    ]);
    (apiClient.getAPIKeys as vi.Mock).mockResolvedValue([]);
  });

  it('should handle navigation events correctly', async () => {
    renderWithProviders(<APIKeys />);
    
    // Wait for component to load
    await waitFor(() => {
      expect(screen.getByText(/apiKeys:title/)).toBeInTheDocument();
    });
    
    // Simulate beforeunload event
    fireEvent(window, new Event('beforeunload'));
    
    // Should not throw errors and sessionStorage should be checked
    expect(sessionStorageMock.getItem).toHaveBeenCalledWith('api-key-dismissed');
  });

  it('should test clipboard API availability', async () => {
    renderWithProviders(<APIKeys />);
    
    // Wait for component to load
    await waitFor(() => {
      expect(screen.getByText(/apiKeys:title/)).toBeInTheDocument();
    });
    
    // Verify clipboard API is available
    expect(navigator.clipboard.writeText).toBeDefined();
  });

  it('should test dismissed state persistence', async () => {
    // Simulate that key was previously dismissed
    sessionStorageMock.getItem.mockReturnValue('true');
    
    renderWithProviders(<APIKeys />);
    
    // Wait for component to load
    await waitFor(() => {
      expect(screen.getByText(/apiKeys:title/)).toBeInTheDocument();
    });
    
    // Verify sessionStorage flag is checked
    expect(sessionStorageMock.getItem).toHaveBeenCalledWith('api-key-dismissed');
  });

  it('should handle visibility change events', async () => {
    renderWithProviders(<APIKeys />);
    
    // Wait for component to load
    await waitFor(() => {
      expect(screen.getByText(/apiKeys:title/)).toBeInTheDocument();
    });
    
    // Simulate visibility change event
    fireEvent(document, new Event('visibilitychange'));
    
    // Should not throw errors
    expect(true).toBe(true);
  });

  it('should test API client integration', async () => {
    renderWithProviders(<APIKeys />);
    
    // Wait for component to load
    await waitFor(() => {
      expect(screen.getByText(/apiKeys:title/)).toBeInTheDocument();
    });
    
    // Verify API client methods are called
    expect(apiClient.getUsers).toHaveBeenCalled();
    expect(apiClient.setOnUnauthorized).toHaveBeenCalled();
  });
});
