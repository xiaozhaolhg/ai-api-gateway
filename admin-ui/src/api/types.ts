// API client interface and types
export interface APIClientInterface {
  // Authentication
  onUnauthorized?: UnauthorizedCallback;
  login(email: string, password: string): Promise<LoginResponse>;
  register(name: string, username: string, email: string, password: string): Promise<RegisterResponse>;
  logout(): Promise<void>;
  getCurrentUser(): Promise<User>;

  // Providers
  getProviders(): Promise<Provider[]>;
  createProvider(provider: Omit<Provider, 'id' | 'created_at' | 'updated_at'>): Promise<Provider>;
  updateProvider(id: string, provider: Partial<Provider>): Promise<Provider>;
  deleteProvider(id: string): Promise<void>;

  // Users
  getUsers(): Promise<User[]>;
  createUser(user: Omit<User, 'id' | 'created_at'>): Promise<User>;
  updateUser(id: string, user: Partial<User>): Promise<User>;
  deleteUser(id: string): Promise<void>;
  checkUsernameAvailability(username: string): Promise<{ available: boolean }>;

  // API Keys
  getAPIKeys(userId: string): Promise<APIKey[]>;
  createAPIKey(userId: string, name: string): Promise<{ api_key_id: string; api_key: string }>;
  deleteAPIKey(id: string): Promise<void>;

  // Usage
  getUsage(userId?: string, startDate?: string, endDate?: string): Promise<UsageRecord[]>;

  // Routing Rules
  getRoutingRules(): Promise<RoutingRule[]>;
  createRoutingRule(rule: Omit<RoutingRule, 'id' | 'created_at' | 'updated_at'>): Promise<RoutingRule>;
  updateRoutingRule(id: string, rule: Partial<RoutingRule>): Promise<RoutingRule>;
  deleteRoutingRule(id: string): Promise<void>;

  // Groups
  getGroups(): Promise<Group[]>;
  createGroup(group: Omit<Group, 'id' | 'created_at' | 'updated_at' | 'member_count'>): Promise<Group>;
  updateGroup(id: string, group: Partial<Group>): Promise<Group>;
  deleteGroup(id: string): Promise<void>;
  addGroupMember(groupId: string, userId: string): Promise<void>;
  removeGroupMember(groupId: string, userId: string): Promise<void>;
  getGroupMembers(groupId: string): Promise<User[]>;
  getGroupPermissions(groupId: string): Promise<Permission[]>;
  getUserGroups(userId: string): Promise<Group[]>;

  // Permissions
  getPermissions(): Promise<Permission[]>;
  createPermission(permission: Omit<Permission, 'id' | 'created_at' | 'updated_at'>): Promise<Permission>;
  updatePermission(id: string, permission: Partial<Permission>): Promise<Permission>;
  deletePermission(id: string): Promise<void>;

  // Budgets
  getBudgets(): Promise<Budget[]>;
  createBudget(budget: Omit<Budget, 'id' | 'created_at' | 'updated_at' | 'current_spend'>): Promise<Budget>;
  updateBudget(id: string, budget: Partial<Budget>): Promise<Budget>;
  deleteBudget(id: string): Promise<void>;

  // Pricing Rules
  getPricingRules(): Promise<PricingRule[]>;
  createPricingRule(rule: Omit<PricingRule, 'id' | 'created_at' | 'updated_at'>): Promise<PricingRule>;
  updatePricingRule(id: string, rule: Partial<PricingRule>): Promise<PricingRule>;
  deletePricingRule(id: string): Promise<void>;

  // Alert Rules
  getAlertRules(): Promise<AlertRule[]>;
  createAlertRule(rule: Omit<AlertRule, 'id' | 'created_at' | 'updated_at'>): Promise<AlertRule>;
  updateAlertRule(id: string, rule: Partial<AlertRule>): Promise<AlertRule>;
  deleteAlertRule(id: string): Promise<void>;

  // Alerts
  getAlerts(): Promise<Alert[]>;
  acknowledgeAlert(id: string): Promise<void>;

  // Health
  getProviderHealth(): Promise<ProviderHealth[]>;
}

// Response Types
export interface LoginResponse {
  token: string;
  user: User;
}

export interface RegisterResponse {
  token: string;
  user: User;
}

// Entity Types (from existing client.ts)
export interface Provider {
  id: string;
  name: string;
  type: string;
  base_url: string;
  models: string[];
  status: string;
  created_at: string;
  updated_at: string;
}

export interface User {
  id: string;
  name: string;
  email: string;
  username: string;
  password?: string;
  role: string;
  status: string;
  groups?: string[];
  created_at: string;
}

export interface APIKey {
  id: string;
  user_id: string;
  name: string;
  scopes: string[];
  created_at: string;
  expires_at?: string;
}

export interface UsageRecord {
  id: string;
  user_id: string;
  model: string;
  provider: string;
  prompt_tokens: number;
  completion_tokens: number;
  total_tokens: number;
  cost: number;
  timestamp: string;
}

export interface RoutingRule {
  id: string;
  model_pattern: string;
  provider: string;
  adapter_type: string;
  priority: number;
  fallback_chain: string[];
  status: string;
  created_at: string;
  updated_at: string;
}

export interface Group {
  id: string;
  name: string;
  description: string;
  member_count: number;
  tier_id?: string;
  parent_group_id?: string;
  created_at: string;
  updated_at: string;
}

export interface Permission {
  id: string;
  group_id: string;
  resource_type: string;
  resource_id: string;
  action: string;
  effect: string;
  created_at: string;
  updated_at: string;
}

export interface Budget {
  id: string;
  name: string;
  scope: string;
  scope_id?: string;
  limit: number;
  current_spend: number;
  period: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface PricingRule {
  id: string;
  model: string;
  provider: string;
  prompt_price: number;
  completion_price: number;
  currency: string;
  effective_date: string;
  created_at: string;
  updated_at: string;
}

export interface AlertRule {
  id: string;
  name: string;
  metric: string;
  condition: string;
  threshold: number;
  channel: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface Alert {
  id: string;
  rule_id: string;
  severity: string;
  status: string;
  triggered_at: string;
  description: string;
  acknowledged_at?: string;
}

export interface ProviderHealth {
  id: string;
  name: string;
  status: string;
  latency_ms: number;
  error_rate: number;
  last_check: string;
}

export interface Tier {
  id: string;
  name: string;
  description: string;
  is_default: boolean;
  allowed_models: string[];
  allowed_providers: string[];
  created_at: string;
  updated_at: string;
}

// Mock Data Store Interface
export interface MockDataStore {
  users: User[];
  providers: Provider[];
  apiKeys: APIKey[];
  usage: UsageRecord[];
  routingRules: RoutingRule[];
  groups: Group[];
  permissions: Permission[];
  budgets: Budget[];
  pricingRules: PricingRule[];
  alertRules: AlertRule[];
  alerts: Alert[];
  providerHealth: ProviderHealth[];
  tiers: Tier[];
}

// API Configuration
export interface APIConfig {
  useMock: boolean;
  mockDelay: number;
  baseURL: string;
}

export type UnauthorizedCallback = () => void;
