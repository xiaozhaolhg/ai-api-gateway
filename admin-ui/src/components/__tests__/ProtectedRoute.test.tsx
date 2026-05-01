import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ProtectedRoute } from '../ProtectedRoute';
import { setMockAuthToken, clearMockAuth } from '../../test/utils';

const renderWithRouter = (component: React.ReactElement) => {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  });
  return render(
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        {component}
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

  it('should redirect to login when not authenticated', () => {
    renderWithRouter(<ProtectedRoute><div data-testid="protected">Protected</div></ProtectedRoute>);
    expect(screen.queryByTestId('protected')).not.toBeInTheDocument();
  });

  it('should render children when authenticated', async () => {
    setMockAuthToken('mock-jwt-token-admin-123');
    renderWithRouter(
      <ProtectedRoute>
        <div data-testid="protected">Protected</div>
      </ProtectedRoute>
    );
    expect(screen.getByTestId('protected')).toBeInTheDocument();
  });

  it('should respect requiredRole for admin', () => {
    setMockAuthToken('mock-jwt-token-admin-123');
    renderWithRouter(
      <ProtectedRoute requiredRole="admin">
        <div data-testid="protected">Admin Only</div>
      </ProtectedRoute>
    );
    expect(screen.getByTestId('protected')).toBeInTheDocument();
  });
});
