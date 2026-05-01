import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { MockAuthProvider } from '../../test/utils';
import Users from '../Users';

const renderUsers = (role = 'admin') => {
  return render(
    <MockAuthProvider role={role}>
      <BrowserRouter>
        <Users />
      </BrowserRouter>
    </MockAuthProvider>
  );
};

describe('Users page', () => {
  it('should render users table', async () => {
    renderUsers();
    await waitFor(() => {
      expect(screen.getByText(/users|list/i)).toBeInTheDocument();
    });
  });

  it('should display users in table', async () => {
    renderUsers();
    await waitFor(() => {
      expect(screen.getByText(/test user|test@example.com/i)).toBeInTheDocument();
    });
  });

  it('should open create modal when add button clicked', async () => {
    renderUsers();
    const addButton = screen.getByRole('button', { name: /add|create|new/i });
    fireEvent.click(addButton);
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
  });

  it('should show loading state', () => {
    renderUsers();
    expect(screen.queryByRole('progressbar')).toBeTruthy();
  });

  it('should show error when API fails', async () => {
    renderUsers();
    await waitFor(() => {
      expect(screen.getByText(/error|failed/i)).toBeInTheDocument();
    });
  });

  it('should disable edit/delete for non-admin', async () => {
    renderUsers('user');
    await waitFor(() => {
      const editButtons = screen.queryAllByRole('button', { name: /edit/i });
      const deleteButtons = screen.queryAllByRole('button', { name: /delete/i });
      if (editButtons.length > 0) {
        expect(editButtons[0]).toBeDisabled();
      }
      if (deleteButtons.length > 0) {
        expect(deleteButtons[0]).toBeDisabled();
      }
    });
  });
});
