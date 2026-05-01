import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { Login } from '../../pages/Login';
import { MockAuthProvider } from '../../test/utils';

const renderLogin = () => {
  return render(
    <MockAuthProvider>
      <Login />
    </MockAuthProvider>
  );
};

describe('Login form validation', () => {
  it('should show error when email is empty', async () => {
    renderLogin();
    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: '' } });
    fireEvent.click(screen.getByRole('button', { name: /login/i }));
    await waitFor(() => {
      expect(screen.getByText(/required|email is required/i)).toBeInTheDocument();
    });
  });

  it('should show error when password is empty', async () => {
    renderLogin();
    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'test@example.com' } });
    fireEvent.click(screen.getByRole('button', { name: /login/i }));
    await waitFor(() => {
      expect(screen.getByText(/password is required/i)).toBeInTheDocument();
    });
  });

  it('should show error for invalid email format', async () => {
    renderLogin();
    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'invalid-email' } });
    fireEvent.change(screen.getByLabelText(/password/i), { target: { value: 'password123' } });
    fireEvent.click(screen.getByRole('button', { name: /login/i }));
    await waitFor(() => {
      expect(screen.getByText(/invalid email|valid email/i)).toBeInTheDocument();
    });
  });

  it('should submit with valid credentials', async () => {
    renderLogin();
    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'test@example.com' } });
    fireEvent.change(screen.getByLabelText(/password/i), { target: { value: 'password123' } });
    fireEvent.click(screen.getByRole('button', { name: /login/i }));
    await waitFor(() => {
      expect(screen.queryByText(/error/i)).not.toBeInTheDocument();
    });
  });
});
