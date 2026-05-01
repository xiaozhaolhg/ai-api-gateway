import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { MockAuthProvider } from '../../test/utils';
import { Budgets } from '../Budgets';

const renderBudgets = (role = 'admin') => {
  return render(
    <MockAuthProvider role={role}>
      <BrowserRouter>
        <Budgets />
      </BrowserRouter>
    </MockAuthProvider>
  );
};

describe('Budgets page', () => {
  it('should render budgets table', async () => {
    renderBudgets();
    await waitFor(() => {
      expect(screen.getByText(/budgets|list/i)).toBeInTheDocument();
    });
  });

  it('should display budgets in table', async () => {
    renderBudgets();
    await waitFor(() => {
      expect(screen.getByText(/test budget|global/i)).toBeInTheDocument();
    });
  });

  it('should open create modal', async () => {
    renderBudgets();
    const addButton = screen.getByRole('button', { name: /add|create|new/i });
    fireEvent.click(addButton);
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
  });

  it('should show budget status badges', async () => {
    renderBudgets();
    await waitFor(() => {
      expect(screen.getByText(/active|warning|exceeded/i)).toBeInTheDocument();
    });
  });

  it('should show loading state', () => {
    renderBudgets();
    expect(screen.queryByRole('progressbar')).toBeTruthy();
  });

  it('should show error when API fails', async () => {
    renderBudgets();
    await waitFor(() => {
      expect(screen.getByText(/error|failed/i)).toBeInTheDocument();
    });
  });

  it('should be hidden for non-admin', async () => {
    renderBudgets('user');
    await waitFor(() => {
      expect(screen.queryByText(/budgets/i)).not.toBeInTheDocument();
    });
  });
});
