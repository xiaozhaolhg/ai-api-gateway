import React from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { GroupMembersTab } from '../../components/GroupMembersTab';
import { apiClient } from '../../api/client';

vi.mock('../../api/client', () => ({
  apiClient: {
    getGroupMembers: vi.fn(),
    getUsers: vi.fn(),
    addGroupMember: vi.fn(),
    removeGroupMember: vi.fn(),
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

describe('GroupMembersTab', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders members table', async () => {
    vi.mocked(apiClient.getGroupMembers).mockResolvedValue([
      { id: '1', name: 'User 1', email: 'user1@example.com', role: 'user' },
    ]);
    vi.mocked(apiClient.getUsers).mockResolvedValue([]);

    renderWithProviders(<GroupMembersTab groupId="group-1" />);

    expect(await screen.findByText('User 1')).toBeTruthy();
    expect(screen.getByText('user1@example.com')).toBeTruthy();
  });

  it('shows add member button', () => {
    vi.mocked(apiClient.getGroupMembers).mockResolvedValue([]);
    vi.mocked(apiClient.getUsers).mockResolvedValue([]);

    renderWithProviders(<GroupMembersTab groupId="group-1" />);

    expect(screen.getByText('Add Member')).toBeTruthy();
  });
});
