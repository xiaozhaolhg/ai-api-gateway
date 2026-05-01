import type {
  User,
  Provider,
  APIKey,
  UsageRecord,
  RoutingRule,
  Group,
  Permission,
  Budget,
  PricingRule,
  AlertRule,
  Alert,
  ProviderHealth,
} from '../api/types';

export const createMockProvider = (overrides: Partial<Provider> = {}): Provider => ({
  id: 'prov_123',
  name: 'Test Provider',
  type: 'ollama',
  base_url: 'http://localhost:11434',
  models: ['llama2', 'mistral'],
  status: 'active',
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
  ...overrides,
});

export const createMockUser = (overrides: Partial<User> = {}): User => ({
  id: 'usr_123',
  name: 'Test User',
  email: 'test@example.com',
  role: 'admin',
  status: 'active',
  created_at: new Date().toISOString(),
  ...overrides,
});

export const createMockAPIKey = (overrides: Partial<APIKey> = {}): APIKey => ({
  id: 'key_123',
  user_id: 'usr_123',
  name: 'Test Key',
  scopes: ['read', 'write'],
  created_at: new Date().toISOString(),
  ...overrides,
});

export const createMockUsageRecord = (overrides: Partial<UsageRecord> = {}): UsageRecord => ({
  id: 'usage_123',
  user_id: 'usr_123',
  model: 'ollama:llama2',
  provider: 'ollama',
  prompt_tokens: 100,
  completion_tokens: 50,
  total_tokens: 150,
  cost: 0.001,
  timestamp: new Date().toISOString(),
  ...overrides,
});

export const createMockRoutingRule = (overrides: Partial<RoutingRule> = {}): RoutingRule => ({
  id: 'rule_123',
  model_pattern: 'ollama:*',
  provider: 'Ollama Local',
  adapter_type: 'ollama',
  priority: 1,
  fallback_chain: [],
  status: 'active',
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
  ...overrides,
});

export const createMockGroup = (overrides: Partial<Group> = {}): Group => ({
  id: 'group_123',
  name: 'Test Group',
  description: 'Test group description',
  member_count: 0,
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
  ...overrides,
});

export const createMockPermission = (overrides: Partial<Permission> = {}): Permission => ({
  id: 'perm_123',
  group_id: 'group_123',
  model_pattern: '*',
  effect: 'allow',
  status: 'active',
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
  ...overrides,
});

export const createMockBudget = (overrides: Partial<Budget> = {}): Budget => ({
  id: 'budget_123',
  name: 'Test Budget',
  scope: 'global',
  limit: 100.0,
  current_spend: 25.5,
  period: 'monthly',
  status: 'active',
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
  ...overrides,
});

export const createMockPricingRule = (overrides: Partial<PricingRule> = {}): PricingRule => ({
  id: 'price_123',
  model: 'ollama:llama2',
  provider: 'ollama',
  prompt_price: 0.001,
  completion_price: 0.002,
  currency: 'USD',
  effective_date: new Date().toISOString(),
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
  ...overrides,
});

export const createMockAlertRule = (overrides: Partial<AlertRule> = {}): AlertRule => ({
  id: 'alert_rule_123',
  name: 'Test Alert Rule',
  metric: 'budget_usage',
  condition: 'greater_than',
  threshold: 80,
  channel: 'email',
  status: 'active',
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
  ...overrides,
});

export const createMockAlert = (overrides: Partial<Alert> = {}): Alert => ({
  id: 'alert_123',
  rule_id: 'alert_rule_123',
  severity: 'warning',
  status: 'active',
  triggered_at: new Date().toISOString(),
  description: 'Test alert description',
  ...overrides,
});

export const createMockProviderHealth = (overrides: Partial<ProviderHealth> = {}): ProviderHealth => ({
  id: 'health_123',
  name: 'Test Provider',
  status: 'healthy',
  latency_ms: 50,
  error_rate: 0.0,
  last_check: new Date().toISOString(),
  ...overrides,
});

export const mockApiResponses = {
  login: {
    token: 'mock-jwt-token',
    user: createMockUser(),
  },
  getProviders: [createMockProvider(), createMockProvider({ id: 'prov_456', name: 'Another Provider' })],
  getUsers: [createMockUser(), createMockUser({ id: 'usr_456', role: 'user' })],
  getAPIKeys: [createMockAPIKey()],
  getUsage: [createMockUsageRecord()],
  getRoutingRules: [createMockRoutingRule()],
  getGroups: [createMockGroup()],
  getPermissions: [createMockPermission()],
  getBudgets: [createMockBudget()],
  getPricingRules: [createMockPricingRule()],
  getAlertRules: [createMockAlertRule()],
  getAlerts: [createMockAlert()],
  getProviderHealth: [createMockProviderHealth()],
};
