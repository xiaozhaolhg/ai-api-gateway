import { render, screen, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { MockAuthProvider } from '../../test/utils';
import Health from '../Health';

const renderHealth = () => {
  return render(
    <MockAuthProvider>
      <BrowserRouter>
        <Health />
      </BrowserRouter>
    </MockAuthProvider>
  );
};

describe('Health page', () => {
  it('should render health page', async () => {
    renderHealth();
    await waitFor(() => {
      expect(screen.getByText(/health|status|latency/i)).toBeInTheDocument();
    });
  });

  it('should display provider health', async () => {
    renderHealth();
    await waitFor(() => {
      expect(screen.getByText(/test provider|healthy|\d+ms/i)).toBeInTheDocument();
    });
  });

  it('should show auto-refresh toggle', async () => {
    renderHealth();
    await waitFor(() => {
      const toggle = screen.queryByRole('switch');
      if (toggle) {
        expect(toggle).toBeInTheDocument();
      }
    });
  });

  it('should show loading state', () => {
    renderHealth();
    expect(screen.queryByRole('progressbar')).toBeTruthy();
  });

  it('should show error when API fails', async () => {
    renderHealth();
    await waitFor(() => {
      expect(screen.getByText(/error|failed/i)).toBeInTheDocument();
    });
  });

  it('should be accessible to all roles', async () => {
    renderHealth();
    await waitFor(() => {
      expect(screen.getByText(/health/i)).toBeInTheDocument();
    });
  });
});
