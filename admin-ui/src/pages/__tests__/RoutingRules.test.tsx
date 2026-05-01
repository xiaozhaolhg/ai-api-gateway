import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { MockAuthProvider } from '../../test/utils';
import { RoutingRules } from '../RoutingRules';

const renderRoutingRules = (role = 'admin') => {
  return render(
    <MockAuthProvider role={role}>
      <BrowserRouter>
        <RoutingRules />
      </BrowserRouter>
    </MockAuthProvider>
  );
};

describe('RoutingRules page', () => {
  it('should render routing rules table', async () => {
    renderRoutingRules();
    await waitFor(() => {
      expect(screen.getByText(/routing rules|rules/i)).toBeInTheDocument();
    });
  });

  it('should display rules in table', async () => {
    renderRoutingRules();
    await waitFor(() => {
      expect(screen.getByText(/ollama:*|pattern/i)).toBeInTheDocument();
    });
  });

  it('should open create modal', async () => {
    renderRoutingRules();
    const addButton = screen.getByRole('button', { name: /add|create|new/i });
    fireEvent.click(addButton);
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
  });

  it('should show loading state', () => {
    renderRoutingRules();
    expect(screen.queryByRole('progressbar')).toBeTruthy();
  });

  it('should show error when API fails', async () => {
    renderRoutingRules();
    await waitFor(() => {
      expect(screen.getByText(/error|failed/i)).toBeInTheDocument();
    });
  });

  it('should be hidden for non-admin', async () => {
    renderRoutingRules('user');
    await waitFor(() => {
      expect(screen.queryByText(/routing rules/i)).not.toBeInTheDocument();
    });
  });
});
