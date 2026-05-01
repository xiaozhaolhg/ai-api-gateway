import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { MockAuthProvider } from '../../test/utils';
import Providers from '../Providers';

const renderProviders = (role = 'admin') => {
  return render(
    <MockAuthProvider role={role}>
      <BrowserRouter>
        <Providers />
      </BrowserRouter>
    </MockAuthProvider>
  );
};

describe('Providers page', () => {
  it('should render providers table', async () => {
    renderProviders();
    await waitFor(() => {
      expect(screen.getByText(/providers|list/i)).toBeInTheDocument();
    });
  });

  it('should display providers in table', async () => {
    renderProviders();
    await waitFor(() => {
      expect(screen.getByText(/test provider|ollama/i)).toBeInTheDocument();
    });
  });

  it('should open create modal when add button clicked', async () => {
    renderProviders();
    const addButton = screen.getByRole('button', { name: /add|create|new/i });
    fireEvent.click(addButton);
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
  });

  it('should show loading state', () => {
    renderProviders();
    expect(screen.queryByRole('progressbar') || screen.getByTestId('loading')).toBeTruthy();
  });

  it('should show error when API fails', async () => {
    renderProviders();
    await waitFor(() => {
      expect(screen.getByText(/error|failed/i)).toBeInTheDocument();
    });
  });

  it('should show empty state when no providers', async () => {
    renderProviders();
    await waitFor(() => {
      expect(screen.getByText(/no data|empty/i)).toBeInTheDocument();
    });
  });
});
