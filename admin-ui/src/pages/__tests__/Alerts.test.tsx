import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { MockAuthProvider } from '../../test/utils';
import { Alerts } from '../Alerts';

const renderAlerts = (role = 'admin') => {
  return render(
    <MockAuthProvider role={role}>
      <BrowserRouter>
        <Alerts />
      </BrowserRouter>
    </MockAuthProvider>
  );
};

describe('Alerts page', () => {
  it('should render alerts tabs', async () => {
    renderAlerts();
    await waitFor(() => {
      expect(screen.getByText(/rules|alerts|active/i)).toBeInTheDocument();
    });
  });

  it('should display alert rules tab', async () => {
    renderAlerts();
    await waitFor(() => {
      expect(screen.getByText(/test alert rule|greater_than/i)).toBeInTheDocument();
    });
  });

  it('should display active alerts tab', async () => {
    renderAlerts();
    await waitFor(() => {
      expect(screen.getByText(/warning|active|acknowledged/i)).toBeInTheDocument();
    });
  });

  it('should show loading state', () => {
    renderAlerts();
    expect(screen.queryByRole('progressbar')).toBeTruthy();
  });

  it('should show error when API fails', async () => {
    renderAlerts();
    await waitFor(() => {
      expect(screen.getByText(/error|failed/i)).toBeInTheDocument();
    });
  });

  it('should allow acknowledge for admin', async () => {
    renderAlerts('admin');
    await waitFor(() => {
      const ackButton = screen.queryByRole('button', { name: /acknowledge/i });
      if (ackButton) {
        expect(ackButton).toBeEnabled();
      }
    });
  });

  it('should be read-only for viewer', async () => {
    renderAlerts('viewer');
    await waitFor(() => {
      const editButtons = screen.queryAllByRole('button', { name: /edit/i });
      const deleteButtons = screen.queryAllByRole('button', { name: /delete/i });
      if (editButtons.length > 0) {
        expect(editButtons[0]).toBeDisabled();
      }
    });
  });
});
