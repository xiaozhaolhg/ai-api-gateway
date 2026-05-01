import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { MockAuthProvider } from '../../test/utils';
import Providers from '../Providers';

const renderProvidersForm = () => {
  return render(
    <MockAuthProvider role="admin">
      <BrowserRouter>
        <Providers />
      </BrowserRouter>
    </MockAuthProvider>
  );
};

describe('Provider form validation', () => {
  it('should show error when name is empty', async () => {
    renderProvidersForm();
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

  it('should show error when base_url is empty', async () => {
    renderProvidersForm();
    const addButton = screen.getByRole('button', { name: /add|create|new/i });
    fireEvent.click(addButton);
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
    fireEvent.click(screen.getByRole('button', { name: /submit|save|confirm/i }));
    await waitFor(() => {
      expect(screen.getByText(/url.*required|base.*required/i)).toBeInTheDocument();
    });
  });

  it('should submit with valid data', async () => {
    renderProvidersForm();
    const addButton = screen.getByRole('button', { name: /add|create|new/i });
    fireEvent.click(addButton);
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
    fireEvent.change(screen.getByLabelText(/name/i), { target: { value: 'New Provider' } });
    fireEvent.change(screen.getByLabelText(/url|base/i), { target: { value: 'http://test.com' } });
    fireEvent.click(screen.getByRole('button', { name: /submit|save|confirm/i }));
    await waitFor(() => {
      expect(screen.queryByText(/error|required/i)).not.toBeInTheDocument();
    });
  });
});
