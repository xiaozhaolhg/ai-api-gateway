import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { MockAuthProvider } from '../../test/utils';
import Users from '../Users';

const renderUsersForm = () => {
  return render(
    <MockAuthProvider role="admin">
      <BrowserRouter>
        <Users />
      </BrowserRouter>
    </MockAuthProvider>
  );
};

describe('User form validation', () => {
  it('should show error when name is empty', async () => {
    renderUsersForm();
    const addButton = screen.getByRole('button', { name: /add|create|new/i });
    fireEvent.click(addButton);
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByRole('button', { name: /submit|save|confirm/i }));
    await waitFor(() => {
      expect(screen.getByText(/name.*required|required field/i)).toBeInTheDocument();
    });
  });

  it('should show error when email is invalid', async () => {
    renderUsersForm();
    const addButton = screen.getByRole('button', { name: /add|create|new/i });
    fireEvent.click(addButton);
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'invalid' } });
    fireEvent.click(screen.getByRole('button', { name: /submit|save|confirm/i }));
    await waitFor(() => {
      expect(screen.getByText(/email.*valid|invalid email/i)).toBeInTheDocument();
    });
  });

  it('should submit with valid data', async () => {
    renderUsersForm();
    const addButton = screen.getByRole('button', { name: /add|create|new/i });
    fireEvent.click(addButton);
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
    fireEvent.change(screen.getByLabelText(/name/i), { target: { value: 'New User' } });
    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'new@example.com' } });
    fireEvent.click(screen.getByRole('button', { name: /submit|save|confirm/i }));
    await waitFor(() => {
      expect(screen.queryByText(/error|required/i)).not.toBeInTheDocument();
    });
  });
});
