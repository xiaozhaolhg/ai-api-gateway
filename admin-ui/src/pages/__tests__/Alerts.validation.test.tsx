import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { MockAuthProvider } from '../../test/utils';
import Alerts from '../Alerts';

const renderAlertRulesForm = () => {
  return render(
    <MockAuthProvider role="admin">
      <BrowserRouter>
        <Alerts />
      </BrowserRouter>
    </MockAuthProvider>
  );
};

describe('AlertRule form validation', () => {
  it('should show error when name is empty', async () => {
    renderAlertRulesForm();
    const tabs = screen.getAllByRole('tab');
    const rulesTab = tabs.find(tab => /rule/i.test(tab.textContent || ''));
    if (rulesTab) fireEvent.click(rulesTab);
    await waitFor(() => {});
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

  it('should show error when threshold is not numeric', async () => {
    renderAlertRulesForm();
    const tabs = screen.getAllByRole('tab');
    const rulesTab = tabs.find(tab => /rule/i.test(tab.textContent || ''));
    if (rulesTab) fireEvent.click(rulesTab);
    await waitFor(() => {});
    const addButton = screen.getByRole('button', { name: /add|create|new/i });
    fireEvent.click(addButton);
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
    fireEvent.change(screen.getByLabelText(/threshold/i), { target: { value: 'abc' } });
    fireEvent.click(screen.getByRole('button', { name: /submit|save|confirm/i }));
    await waitFor(() => {
      expect(screen.getByText(/numeric|number|invalid/i)).toBeInTheDocument();
    });
  });

  it('should submit with valid data', async () => {
    renderAlertRulesForm();
    const tabs = screen.getAllByRole('tab');
    const rulesTab = tabs.find(tab => /rule/i.test(tab.textContent || ''));
    if (rulesTab) fireEvent.click(rulesTab);
    await waitFor(() => {});
    const addButton = screen.getByRole('button', { name: /add|create|new/i });
    fireEvent.click(addButton);
    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
    fireEvent.change(screen.getByLabelText(/name/i), { target: { value: 'New Rule' } });
    fireEvent.change(screen.getByLabelText(/threshold/i), { target: { value: '80' } });
    fireEvent.click(screen.getByRole('button', { name: /submit|save|confirm/i }));
    await waitFor(() => {
      expect(screen.queryByText(/error|required/i)).not.toBeInTheDocument();
    });
  });
});
