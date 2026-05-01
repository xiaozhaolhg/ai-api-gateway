import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import MockAPIClient from '../mockClient';
import MockDataHandler from '../../mock/handlers/dataHandler';

describe('MockAPIClient', () => {
  let mockClient: MockAPIClient;
  let dataHandler: MockDataHandler;

  beforeEach(() => {
    // Create a fresh instance for each test with no delay
    mockClient = new MockAPIClient(0);
    dataHandler = MockDataHandler.getInstance();
    // Reset to defaults before each test
    dataHandler.resetToDefaults();
  });

  afterEach(() => {
    // Clean up after each test
    dataHandler.resetToDefaults();
  });

  describe('Authentication', () => {
    it('should login with valid email', async () => {
      const result = await mockClient.login('admin@example.com', 'password');
      expect(result).toHaveProperty('token');
      expect(result).toHaveProperty('user');
      expect(result.user.email).toBe('admin@example.com');
    });

    it('should login with valid username', async () => {
      const result = await mockClient.login('testuser', 'password');
      expect(result).toHaveProperty('token');
      expect(result).toHaveProperty('user');
      expect(result.user.email).toBe('testuser@local.dev');
    });

    it('should fail login with non-existent user', async () => {
      await expect(
        mockClient.login('nonexistent@example.com', 'password')
      ).rejects.toThrow('User not found');
    });

    it('should fail login with invalid password', async () => {
      await expect(
        mockClient.login('admin@example.com', '')
      ).rejects.toThrow('Invalid password');
    });

    it('should register a new user', async () => {
      const result = await mockClient.register('New User', 'newuser', 'newuser@example.com', 'password123');
      expect(result).toHaveProperty('token');
      expect(result.user.name).toBe('New User');
      expect(result.user.email).toBe('newuser@example.com');
    });

    it('should fail registration with duplicate email', async () => {
      await expect(
        mockClient.register('Duplicate', 'duplicate', 'admin@example.com', 'password')
      ).rejects.toThrow('User already exists');
    });

    it('should logout successfully', async () => {
      await expect(mockClient.logout()).resolves.toBeUndefined();
    });

    it('should get current user', async () => {
      const user = await mockClient.getCurrentUser();
      expect(user).toHaveProperty('id');
      expect(user).toHaveProperty('name');
      expect(user).toHaveProperty('email');
    });
  });

  describe('Providers', () => {
    it('should get all providers', async () => {
      const providers = await mockClient.getProviders();
      expect(Array.isArray(providers)).toBe(true);
      expect(providers.length).toBeGreaterThan(0);
    });

    it('should create a new provider', async () => {
      const newProvider = {
        name: 'Test Provider',
        type: 'test',
        base_url: 'http://test.com',
        models: ['model1', 'model2'],
        status: 'active'
      };
      
      const created = await mockClient.createProvider(newProvider);
      expect(created).toHaveProperty('id');
      expect(created.name).toBe('Test Provider');
    });

    it('should update a provider', async () => {
      const providers = await mockClient.getProviders();
      const firstProvider = providers[0];
      
      const updated = await mockClient.updateProvider(firstProvider.id, {
        name: 'Updated Name'
      });
      
      expect(updated.name).toBe('Updated Name');
    });

    it('should delete a provider', async () => {
      const providers = await mockClient.getProviders();
      const firstProvider = providers[0];
      
      await expect(mockClient.deleteProvider(firstProvider.id)).resolves.toBeUndefined();
      
      const afterDelete = await mockClient.getProviders();
      expect(afterDelete.find(p => p.id === firstProvider.id)).toBeUndefined();
    });
  });

  describe('Users', () => {
    it('should get all users', async () => {
      const users = await mockClient.getUsers();
      expect(Array.isArray(users)).toBe(true);
      expect(users.length).toBeGreaterThan(0);
    });

    it('should create a new user', async () => {
      const newUser = {
        name: 'Test User',
        email: 'test@example.com',
        role: 'user',
        status: 'active'
      };
      
      const created = await mockClient.createUser(newUser);
      expect(created).toHaveProperty('id');
      expect(created.name).toBe('Test User');
    });

    it('should update a user', async () => {
      const users = await mockClient.getUsers();
      const firstUser = users[0];
      
      const updated = await mockClient.updateUser(firstUser.id, {
        name: 'Updated Name'
      });
      
      expect(updated.name).toBe('Updated Name');
    });

    it('should delete a user', async () => {
      const users = await mockClient.getUsers();
      const firstUser = users[0];
      
      await expect(mockClient.deleteUser(firstUser.id)).resolves.toBeUndefined();
      
      const afterDelete = await mockClient.getUsers();
      expect(afterDelete.find(u => u.id === firstUser.id)).toBeUndefined();
    });
  });

  describe('Routing Rules', () => {
    it('should get all routing rules', async () => {
      const rules = await mockClient.getRoutingRules();
      expect(Array.isArray(rules)).toBe(true);
      expect(rules.length).toBeGreaterThan(0);
    });

    it('should create a new routing rule', async () => {
      const newRule = {
        model_pattern: 'test:*',
        provider: 'Test Provider',
        adapter_type: 'test',
        priority: 1,
        fallback_chain: [],
        status: 'active'
      };
      
      const created = await mockClient.createRoutingRule(newRule);
      expect(created).toHaveProperty('id');
      expect(created.model_pattern).toBe('test:*');
    });

    it('should update a routing rule', async () => {
      const rules = await mockClient.getRoutingRules();
      const firstRule = rules[0];
      
      const updated = await mockClient.updateRoutingRule(firstRule.id, {
        model_pattern: 'updated:*'
      });
      
      expect(updated.model_pattern).toBe('updated:*');
    });

    it('should delete a routing rule', async () => {
      const rules = await mockClient.getRoutingRules();
      const firstRule = rules[0];
      
      await expect(mockClient.deleteRoutingRule(firstRule.id)).resolves.toBeUndefined();
      
      const afterDelete = await mockClient.getRoutingRules();
      expect(afterDelete.find(r => r.id === firstRule.id)).toBeUndefined();
    });
  });

  describe('Budgets', () => {
    it('should get all budgets', async () => {
      const budgets = await mockClient.getBudgets();
      expect(Array.isArray(budgets)).toBe(true);
      expect(budgets.length).toBeGreaterThan(0);
    });

    it('should create a new budget', async () => {
      const newBudget = {
        name: 'Test Budget',
        scope: 'global',
        limit: 100.0,
        period: 'monthly',
        status: 'active'
      };
      
      const created = await mockClient.createBudget(newBudget);
      expect(created).toHaveProperty('id');
      expect(created.name).toBe('Test Budget');
      expect(created.current_spend).toBe(0);
    });

    it('should update a budget', async () => {
      const budgets = await mockClient.getBudgets();
      const firstBudget = budgets[0];
      
      const updated = await mockClient.updateBudget(firstBudget.id, {
        limit: 200.0
      });
      
      expect(updated.limit).toBe(200.0);
    });

    it('should delete a budget', async () => {
      const budgets = await mockClient.getBudgets();
      const firstBudget = budgets[0];
      
      await expect(mockClient.deleteBudget(firstBudget.id)).resolves.toBeUndefined();
      
      const afterDelete = await mockClient.getBudgets();
      expect(afterDelete.find(b => b.id === firstBudget.id)).toBeUndefined();
    });
  });

  describe('Pricing Rules', () => {
    it('should get all pricing rules', async () => {
      const rules = await mockClient.getPricingRules();
      expect(Array.isArray(rules)).toBe(true);
      expect(rules.length).toBeGreaterThan(0);
    });

    it('should create a new pricing rule', async () => {
      const newRule = {
        model: 'test:model',
        provider: 'test',
        prompt_price: 0.001,
        completion_price: 0.002,
        currency: 'USD',
        effective_date: '2024-01-01T00:00:00Z'
      };
      
      const created = await mockClient.createPricingRule(newRule);
      expect(created).toHaveProperty('id');
      expect(created.model).toBe('test:model');
    });

    it('should update a pricing rule', async () => {
      const rules = await mockClient.getPricingRules();
      const firstRule = rules[0];
      
      const updated = await mockClient.updatePricingRule(firstRule.id, {
        prompt_price: 0.005
      });
      
      expect(updated.prompt_price).toBe(0.005);
    });

    it('should delete a pricing rule', async () => {
      const rules = await mockClient.getPricingRules();
      const firstRule = rules[0];
      
      await expect(mockClient.deletePricingRule(firstRule.id)).resolves.toBeUndefined();
      
      const afterDelete = await mockClient.getPricingRules();
      expect(afterDelete.find(r => r.id === firstRule.id)).toBeUndefined();
    });
  });

  describe('Alert Rules', () => {
    it('should get all alert rules', async () => {
      const rules = await mockClient.getAlertRules();
      expect(Array.isArray(rules)).toBe(true);
      expect(rules.length).toBeGreaterThan(0);
    });

    it('should create a new alert rule', async () => {
      const newRule = {
        name: 'Test Alert',
        metric: 'test_metric',
        condition: 'greater_than',
        threshold: 100,
        channel: 'email',
        status: 'active'
      };
      
      const created = await mockClient.createAlertRule(newRule);
      expect(created).toHaveProperty('id');
      expect(created.name).toBe('Test Alert');
    });

    it('should update an alert rule', async () => {
      const rules = await mockClient.getAlertRules();
      const firstRule = rules[0];
      
      const updated = await mockClient.updateAlertRule(firstRule.id, {
        threshold: 200
      });
      
      expect(updated.threshold).toBe(200);
    });

    it('should delete an alert rule', async () => {
      const rules = await mockClient.getAlertRules();
      const firstRule = rules[0];
      
      await expect(mockClient.deleteAlertRule(firstRule.id)).resolves.toBeUndefined();
      
      const afterDelete = await mockClient.getAlertRules();
      expect(afterDelete.find(r => r.id === firstRule.id)).toBeUndefined();
    });
  });

  describe('Alerts', () => {
    it('should get all alerts', async () => {
      const alerts = await mockClient.getAlerts();
      expect(Array.isArray(alerts)).toBe(true);
    });

    it('should acknowledge an alert', async () => {
      const alerts = await mockClient.getAlerts();
      if (alerts.length > 0) {
        const firstAlert = alerts[0];
        
        await expect(mockClient.acknowledgeAlert(firstAlert.id)).resolves.toBeUndefined();
        
        const afterAcknowledge = await mockClient.getAlerts();
        const acknowledged = afterAcknowledge.find(a => a.id === firstAlert.id);
        expect(acknowledged?.status).toBe('acknowledged');
        expect(acknowledged?.acknowledged_at).toBeDefined();
      }
    });
  });

  describe('Network Delay Simulation', () => {
    it('should simulate network delay', async () => {
      const delayedClient = new MockAPIClient(100);
      const start = Date.now();
      await delayedClient.getProviders();
      const end = Date.now();
      
      expect(end - start).toBeGreaterThanOrEqual(90); // Allow some timing variance
    });
  });
});
