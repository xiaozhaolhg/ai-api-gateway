import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { MockAuthProvider } from '../../test/utils';
import { Groups } from '../Groups';

const renderGroups = (role = 'admin') => {
  return render(
    <MockAuthProvider role={role}>
      <BrowserRouter>
        <Groups />
      </BrowserRouter>
    </MockAuthProvider>
  );
};

describe('Groups page', () => {
  it('should render groups table', async () => {
    renderGroups();
    await waitFor(() => {
      expect(screen.getByText(/groups|list/i)).toBeInTheDocument();
    });
  });

  it('should display groups in table', async () => {
    renderGroups();
    await waitFor(() => {
      expect(screen.getByText(/test group/i)).toBeInTheDocument();
    });
  });

  it('should open create modal', async () => {
    renderGroups();
    const addButton = screen.getByRole('button', { name: /add|create|new/i });
    fireEvent.click(addButton);
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
  });

  it('should show loading state', () => {
    renderGroups();
    expect(screen.queryByRole('progressbar')).toBeTruthy();
  });

  it('should show error when API fails', async () => {
    renderGroups();
    await waitFor(() => {
      expect(screen.getByText(/error|failed/i)).toBeInTheDocument();
    });
  });

  it('should be hidden for non-admin', async () => {
    renderGroups('user');
    await waitFor(() => {
      expect(screen.queryByText(/groups/i)).not.toBeInTheDocument();
    });
  });
});
