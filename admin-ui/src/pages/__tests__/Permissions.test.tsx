import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { MockAuthProvider } from '../../test/utils';
import { Permissions } from '../Permissions';

const renderPermissions = (role = 'admin') => {
  return render(
    <MockAuthProvider role={role}>
      <BrowserRouter>
        <Permissions />
      </BrowserRouter>
    </MockAuthProvider>
  );
};

describe('Permissions page', () => {
  it('should render permissions table', async () => {
    renderPermissions();
    await waitFor(() => {
      expect(screen.getByText(/permissions|list/i)).toBeInTheDocument();
    });
  });

  it('should display permissions in table', async () => {
    renderPermissions();
    await waitFor(() => {
      expect(screen.getByText(/\*|allow/i)).toBeInTheDocument();
    });
  });

  it('should open create modal', async () => {
    renderPermissions();
    const addButton = screen.getByRole('button', { name: /add|create|new/i });
    fireEvent.click(addButton);
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
  });

  it('should show loading state', () => {
    renderPermissions();
    expect(screen.queryByRole('progressbar')).toBeTruthy();
  });

  it('should show error when API fails', async () => {
    renderPermissions();
    await waitFor(() => {
      expect(screen.getByText(/error|failed/i)).toBeInTheDocument();
    });
  });

  it('should be hidden for non-admin', async () => {
    renderPermissions('user');
    await waitFor(() => {
      expect(screen.queryByText(/permissions/i)).not.toBeInTheDocument();
    });
  });
});
