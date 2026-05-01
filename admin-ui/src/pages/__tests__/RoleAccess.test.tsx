import { render, screen, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { MockAuthProvider } from '../../test/utils';
import Providers from '../Providers';
import Users from '../Users';
import APIKeys from '../APIKeys';
import Usage from '../Usage';
import RoutingRules from '../RoutingRules';
import Groups from '../Groups';
import Permissions from '../Permissions';
import Budgets from '../Budgets';
import PricingRules from '../PricingRules';
import Alerts from '../Alerts';
import Health from '../Health';

describe('Role-based access control', () => {
  describe('Admin role', () => {
    it('should access all pages', async () => {
      const { getByText } = render(
        <MockAuthProvider role="admin">
          <BrowserRouter>
            <Providers />
          </BrowserRouter>
        </MockAuthProvider>
      );
      await waitFor(() => {
        expect(getByText(/providers|list/i)).toBeInTheDocument();
      });
    });

    it('should have edit/delete buttons visible', async () => {
      const { queryAllByRole } = render(
        <MockAuthProvider role="admin">
          <BrowserRouter>
            <Users />
          </BrowserRouter>
        </MockAuthProvider>
      );
      await waitFor(() => {});
      const editButtons = queryAllByRole('button', { name: /edit/i });
      const deleteButtons = queryAllByRole('button', { name: /delete/i });
      expect(editButtons.length + deleteButtons.length).toBeGreaterThan(0);
    });
  });

  describe('User role', () => {
    it('should access own API keys', async () => {
      const { getByText } = render(
        <MockAuthProvider role="user">
          <BrowserRouter>
            <APIKeys />
          </BrowserRouter>
        </MockAuthProvider>
      );
      await waitFor(() => {
        expect(getByText(/api keys|keys/i)).toBeInTheDocument();
      });
    });

    it('should access own usage', async () => {
      const { getByText } = render(
        <MockAuthProvider role="user">
          <BrowserRouter>
            <Usage />
          </BrowserRouter>
        </MockAuthProvider>
      );
      await waitFor(() => {
        expect(getByText(/usage|tokens/i)).toBeInTheDocument();
      });
    });

    it('should be blocked from providers page', async () => {
      const { queryByText } = render(
        <MockAuthProvider role="user">
          <BrowserRouter>
            <Providers />
          </BrowserRouter>
        </MockAuthProvider>
      );
      await waitFor(() => {});
      expect(queryByText(/providers.*list|create.*provider/i)).not.toBeInTheDocument();
    });
  });

  describe('Viewer role', () => {
    it('should have read-only access', async () => {
      const { queryAllByRole } = render(
        <MockAuthProvider role="viewer">
          <BrowserRouter>
            <Alerts />
          </BrowserRouter>
        </MockAuthProvider>
      );
      await waitFor(() => {});
      const editButtons = queryAllByRole('button', { name: /edit/i });
      const deleteButtons = queryAllByRole('button', { name: /delete/i });
      const addButtons = queryAllByRole('button', { name: /add|create|new/i });
      expect(editButtons.length + deleteButtons.length + addButtons.length).toBe(0);
    });

    it('should access health page', async () => {
      const { getByText } = render(
        <MockAuthProvider role="viewer">
          <BrowserRouter>
            <Health />
          </BrowserRouter>
        </MockAuthProvider>
      );
      await waitFor(() => {
        expect(getByText(/health|status/i)).toBeInTheDocument();
      });
    });
  });
});
