import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { MockAuthProvider } from '../../test/utils';
import { PricingRules } from '../PricingRules';

const renderPricingRules = (role = 'admin') => {
  return render(
    <MockAuthProvider role={role}>
      <BrowserRouter>
        <PricingRules />
      </BrowserRouter>
    </MockAuthProvider>
  );
};

describe('PricingRules page', () => {
  it('should render pricing rules table', async () => {
    renderPricingRules();
    await waitFor(() => {
      expect(screen.getByText(/pricing rules|pricing/i)).toBeInTheDocument();
    });
  });

  it('should display rules in table', async () => {
    renderPricingRules();
    await waitFor(() => {
      expect(screen.getByText(/0.001|ollama/i)).toBeInTheDocument();
    });
  });

  it('should open create modal', async () => {
    renderPricingRules();
    const addButton = screen.getByRole('button', { name: /add|create|new/i });
    fireEvent.click(addButton);
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
  });

  it('should show loading state', () => {
    renderPricingRules();
    expect(screen.queryByRole('progressbar')).toBeTruthy();
  });

  it('should show error when API fails', async () => {
    renderPricingRules();
    await waitFor(() => {
      expect(screen.getByText(/error|failed/i)).toBeInTheDocument();
    });
  });

  it('should be hidden for non-admin', async () => {
    renderPricingRules('user');
    await waitFor(() => {
      expect(screen.queryByText(/pricing rules/i)).not.toBeInTheDocument();
    });
  });
});
