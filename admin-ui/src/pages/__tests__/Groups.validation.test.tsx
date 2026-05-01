import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { MockAuthProvider } from '../../test/utils';
import Groups from '../Groups';

const renderGroupsForm = () => {
  return render(
    <MockAuthProvider role="admin">
      <BrowserRouter>
        <Groups />
      </BrowserRouter>
    </MockAuthProvider>
  );
};

describe('Group form validation', () => {
  it('should show error when name is empty', async () => {
    renderGroupsForm();
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

  it('should submit with valid data', async () => {
    renderGroupsForm();
    const addButton = screen.getByRole('button', { name: /add|create|new/i });
    fireEvent.click(addButton);
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
    fireEvent.change(screen.getByLabelText(/name/i), { target: { value: 'New Group' } });
    fireEvent.change(screen.getByLabelText(/description/i), { target: { value: 'Group description' } });
    fireEvent.click(screen.getByRole('button', { name: /submit|save|confirm/i }));
    await waitFor(() => {
      expect(screen.queryByText(/error|required/i)).not.toBeInTheDocument();
    });
  });
});
