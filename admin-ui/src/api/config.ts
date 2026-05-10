import type { APIConfig } from './types';

export const API_CONFIG: APIConfig = {
  useMock: import.meta.env.VITE_USE_MOCK === 'true',
  mockDelay: parseInt(import.meta.env.VITE_MOCK_DELAY || '500', 10),
  baseURL: import.meta.env.VITE_API_BASE_URL || '',
};

// Helper function to check if we're in development mode
export const isDevelopment = import.meta.env.DEV;

// Helper function to get current API mode
export const getAPIMode = () => (API_CONFIG.useMock ? 'Mock' : 'Real');
