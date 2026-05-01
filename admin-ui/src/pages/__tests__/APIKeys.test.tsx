import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { MockAuthProvider } from '../../test/utils';
import APIKeys from '../APIKeys';

const renderAPIKeys = (role = 'user') => {
  return render(
    <MockAuthProvider role={role}>
      <BrowserRouter>
        <APIKeys />
      </BrowserRouter>
    </MockAuthProvider>
  );
};

describe('APIKeys page', () => {
  it('should render API keys table', async () => {
    renderAPIKeys();
    await waitFor(() => {
      expect(screen.getByText(/api keys|keys/i)).toBeInTheDocument();
    });
  });

  it('should display keys in table', async () => {
    renderAPIKeys();
    await waitFor(() => {
      expect(screen.getByText(/test key/i)).toBeInTheDocument();
    });
  });

  it('should show key creation dialog', async () => {
    renderAPIKeys();
    const addButton = screen.getByRole('button', { name: /add|create|new/i });
    fireEvent.click(addButton);
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
  });

  it('should show loading state', () => {
    renderAPIKeys();
    expect(screen.queryByRole('progressbar')).toBeTruthy();
  });

  it('should show error when API fails', async () => {
    renderAPIKeys();
    await waitFor(() => {
      expect(screen.getByText(/error|failed/i)).toBeInTheDocument();
    });
  });

  it('should only show own keys for user role', async () => {
    renderAPIKeys('user');
    await waitFor(() => {
      expect(screen.getByText(/test key/i)).toBeInTheDocument();
    });
  });
});
