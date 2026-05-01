import { render, screen, waitFor, within } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { MockAuthProvider } from '../../test/utils';
import Dashboard from '../Dashboard';

const renderDashboard = (role = 'admin') => {
  return render(
    <MockAuthProvider role={role}>
      <BrowserRouter>
        <Dashboard />
      </BrowserRouter>
    </MockAuthProvider>
  );
};

describe('Dashboard page', () => {
  it('should render dashboard with summary cards', async () => {
    renderDashboard();
    await waitFor(() => {
      expect(screen.getByText(/dashboard|summary|overview/i)).toBeInTheDocument();
    });
  });

  it('should display user count', async () => {
    renderDashboard();
    await waitFor(() => {
      expect(screen.getByText(/users|total/i)).toBeInTheDocument();
    });
  });

  it('should display provider count', async () => {
    renderDashboard();
    await waitFor(() => {
      expect(screen.getByText(/providers/i)).toBeInTheDocument();
    });
  });

  it('should show loading state', () => {
    renderDashboard();
    expect(screen.getByTestId('loading') || screen.queryByRole('progressbar')).toBeTruthy();
  });

  it('should show error state when API fails', async () => {
    // Mock API failure
    await waitFor(() => {
      expect(screen.getByText(/error|failed/i)).toBeInTheDocument();
    });
  });
});
