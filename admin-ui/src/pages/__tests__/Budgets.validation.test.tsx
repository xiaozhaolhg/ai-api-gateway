import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { MockAuthProvider } from '../../test/utils';
import Budgets from '../Budgets';

const renderBudgetsForm = () => {
  return render(
    <MockAuthProvider role="admin">
      <BrowserRouter>
        <Budgets />
      </BrowserRouter>
    </MockAuthProvider>
  );
};

describe('Budget form validation', () => {
  it('should show error when name is empty', async () => {
    renderBudgetsForm();
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

  it('should show error when limit is not numeric', async () => {
    renderBudgetsForm();
    const addButton = screen.getByRole('button', { name: /add|create|new/i });
    fireEvent.click(addButton);
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
    fireEvent.change(screen.getByLabelText(/limit/i), { target: { value: 'abc' } });
    fireEvent.click(screen.getByRole('button', { name: /submit|save|confirm/i }));
    await waitFor(() => {
      expect(screen.getByText(/numeric|number|invalid/i)).toBeInTheDocument();
    });
  });

  it('should submit with valid data', async () => {
    renderBudgetsForm();
    const addButton = screen.getByRole('button', { name: /add|create|new/i });
    fireEvent.click(addButton);
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
    fireEvent.change(screen.getByLabelText(/name/i), { target: { value: 'New Budget' } });
    fireEvent.change(screen.getByLabelText(/limit/i), { target: { value: '100' } });
    fireEvent.click(screen.getByRole('button', { name: /submit|save|confirm/i }));
    await waitFor(() => {
      expect(screen.queryByText(/error|required/i)).not.toBeInTheDocument();
    });
  });
});
