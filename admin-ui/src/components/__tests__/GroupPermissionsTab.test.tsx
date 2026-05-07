import React from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { GroupPermissionsTab } from '../../components/GroupPermissionsTab';
import { apiClient } from '../../api/client';

vi.mock('../../api/client', () => ({
  apiClient: {
    getGroupPermissions: vi.fn(),
    createPermission: vi.fn(),
    deletePermission: vi.fn(),
  },
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string) => key }),
}));

Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation((query) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
});

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
      {component}
    </QueryClientProvider>
  );
};

describe('GroupPermissionsTab', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders permissions table', async () => {
    vi.mocked(apiClient.getGroupPermissions).mockResolvedValue([
      { id: '1', group_id: 'group-1', resource_type: 'model', resource_id: 'gpt-4', action: 'access', effect: 'allow' },
    ]);

    renderWithProviders(<GroupPermissionsTab groupId="group-1" />);

    expect(await screen.findByText('model')).toBeTruthy();
    expect(screen.getByText('gpt-4')).toBeTruthy();
  });

  it('shows add permission button', () => {
    vi.mocked(apiClient.getGroupPermissions).mockResolvedValue([]);

    renderWithProviders(<GroupPermissionsTab groupId="group-1" />);

    expect(screen.getByText('Add Permission')).toBeTruthy();
  });
});
